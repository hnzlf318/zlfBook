<template>
    <f7-page class="page-navbar-bottom" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title>{{ tt('OCR Bill Recognition') }}</f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link v-if="recognizedList.length"
                         icon-f7="photo"
                         :title="tt('Select Another Image')"
                         @click="reset">
                </f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block v-if="!recognizedList.length" class="ocr-select-block">
            <p class="ocr-hint">{{ tt('You can select a bill or transaction list screenshot to recognize.') }}</p>
            <f7-button large fill @click="triggerFileInput">{{ recognizing ? tt('Recognizing...') : tt('Select Image') }}</f7-button>
            <input ref="imageInputRef" type="file" class="display-none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="onFileChange" />
        </f7-block>

        <f7-popup v-model:opened="previewOpened">
            <f7-page class="page-navbar-bottom ocr-preview-page">
                <f7-navbar class="ocr-preview-navbar">
                    <f7-nav-left>
                        <f7-link icon-f7="chevron_left" @click="onPreviewCancel"></f7-link>
                    </f7-nav-left>
                    <f7-nav-title>{{ tt('Recognized Image') }}</f7-nav-title>
                    <f7-nav-right class="navbar-compact-icons">
                        <f7-link icon-f7="magnifyingglass"
                                :disabled="recognizing || !previewFile"
                                :title="recognizing ? tt('Recognizing...') : tt('Recognize')"
                                @click="onPreviewRecognize">
                        </f7-link>
                        <f7-link icon-f7="photo"
                                :title="tt('Select Image')"
                                @click="onPreviewSelectImage">
                        </f7-link>
                    </f7-nav-right>
                </f7-navbar>
                <f7-block class="ocr-preview-block">
                    <img v-if="previewImageSrc" :src="previewImageSrc" class="ocr-preview-img" />
                </f7-block>
            </f7-page>
        </f7-popup>

        <f7-block v-if="recognizedList.length">
            <p class="ocr-hint">{{ tt('Recognized transactions (click Add to confirm and edit):') }}</p>
            <div class="list strong inset dividers">
                <div v-for="(item, idx) in recognizedList" :key="idx" class="ocr-result-item">
                    <div class="ocr-fields">
                        <div class="ocr-field"><span class="ocr-label">{{ tt('Income or Expense') }}：</span><span class="ocr-value">{{ getTypeLabel(item.type) }}</span></div>
                        <div class="ocr-field"><span class="ocr-label">{{ tt('Amount') }}：</span><span class="ocr-value">{{ formatAmount(item.sourceAmount) }}</span></div>
                        <div class="ocr-field"><span class="ocr-label">{{ tt('Date') }}：</span><span class="ocr-value">{{ formatTime(item.time) }}</span></div>
                        <div class="ocr-field" v-if="!hiddenFields.category"><span class="ocr-label">{{ tt('Category') }}：</span><span class="ocr-value">{{ getCategoryName(item.categoryId) || '-' }}</span></div>
                        <div class="ocr-field"><span class="ocr-label">{{ tt('Account') }}：</span><span class="ocr-value">{{ getAccountDisplay(item) }}</span></div>
                        <div class="ocr-field" v-if="!hiddenFields.items && getItemNames(item.itemIds)"><span class="ocr-label">{{ tt('Transaction Items') }}：</span><span class="ocr-value">{{ getItemNames(item.itemIds) }}</span></div>
                        <div class="ocr-field" v-if="!hiddenFields.tags && getTagNames(item.tagIds)"><span class="ocr-label">{{ tt('Tags') }}：</span><span class="ocr-value">{{ getTagNames(item.tagIds) }}</span></div>
                        <div class="ocr-field" v-if="item.comment"><span class="ocr-label">{{ tt('Description') }}：</span><span class="ocr-value ocr-desc">{{ item.comment }}</span></div>
                    </div>
                    <div class="ocr-add-wrap">
                        <f7-button small fill color="blue" class="ocr-add-btn" :disabled="addedRowIndices.has(idx)" @click="onAdd(item, idx)">
                            {{ tt('Add') }}
                        </f7-button>
                    </div>
                </div>
            </div>
        </f7-block>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionItemsStore } from '@/stores/transactionItem.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { TransactionType } from '@/core/transaction.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import type { RecognizedReceiptImageResponse } from '@/models/large_language_model.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionItemsStore = useTransactionItemsStore();
const transactionsStore = useTransactionsStore();

const imageInputRef = useTemplateRef<HTMLInputElement>('imageInputRef');

const recognizing = ref(false);
const previewOpened = ref(false);
const previewImageSrc = ref<string | null>(null);
const previewFile = ref<File | null>(null);
const recognizedList = ref<RecognizedReceiptImageResponse[]>([]);
const addedRowIndices = ref<Set<number>>(new Set());

// 可配置的隐藏字段：从 OCR 识别响应中读取（仅在首次识别时读取一次）
const hiddenFields = ref({
    category: true,      // 分类
    items: true,          // 交易项目
    tags: true,           // 标签
});

// 从 OCR 识别响应中提取配置（仅在首次识别时应用）
function applyOCRConfig(config?: { hideCategoryColumn?: boolean, hideItemsColumn?: boolean, hideTagsColumn?: boolean }): void {
    if (!config) return;
    
    if (config.hideCategoryColumn !== undefined) {
        hiddenFields.value.category = config.hideCategoryColumn;
    }
    if (config.hideItemsColumn !== undefined) {
        hiddenFields.value.items = config.hideItemsColumn;
    }
    if (config.hideTagsColumn !== undefined) {
        hiddenFields.value.tags = config.hideTagsColumn;
    }
}

function getCategoryName(categoryId?: string): string {
    if (!categoryId) return '-';
    const cat = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];
    return cat?.name ?? '-';
}

function getAccountName(accountId?: string): string {
    if (!accountId) return '-';
    const acc = accountsStore.allAccountsMap[accountId];
    return acc?.name ?? '-';
}

function getAccountDisplay(item: RecognizedReceiptImageResponse): string {
    if (item.account && item.account.trim().length > 0) {
        return item.account;
    }
    return getAccountName(item.sourceAccountId) || '-';
}

function getTagNames(tagIds?: string[]): string {
    if (!tagIds?.length) return '';
    return tagIds.map(id => transactionTagsStore.allTransactionTagsMap[id]?.name ?? id).join(', ');
}

function getItemNames(itemIds?: string[]): string {
    if (!itemIds?.length) return '';
    const itemsMap = transactionItemsStore.allTransactionItemsMap;
    return itemIds.map(id => itemsMap[id]?.name ?? id).join(', ');
}

function getTypeLabel(type?: number): string {
    if (type === TransactionType.Income) return tt('Income');
    if (type === TransactionType.Expense) return tt('Expense');
    return '-';
}

function formatAmount(amount?: number): string {
    if (amount === undefined || amount === null) return '-';
    const n = amount / 100;
    return n >= 0 ? `¥${n.toFixed(2)}` : `-¥${(-n).toFixed(2)}`;
}

function formatTime(time?: number): string {
    if (time === undefined || time === null) return '-';
    const d = new Date(time * 1000);
    const year = d.getFullYear();
    const month = d.getMonth() + 1;
    const day = d.getDate();
    const hour = d.getHours();
    const minute = d.getMinutes();
    const pad = (n: number) => n.toString().padStart(2, '0');
    return `${year}-${pad(month)}-${pad(day)} ${pad(hour)}:${pad(minute)}`;
}

function triggerFileInput(): void {
    if (recognizing.value) return;
    imageInputRef.value?.click();
}

function onFileChange(event: Event): void {
    const el = event?.target as HTMLInputElement;
    if (!el?.files?.length || !el.files[0]) return;
    const file = el.files[0] as File;
    el.value = '';
    if (previewImageSrc.value) {
        URL.revokeObjectURL(previewImageSrc.value);
    }
    previewFile.value = file;
    previewImageSrc.value = URL.createObjectURL(file);
    previewOpened.value = true;
}

function recognizeFile(file: File): void {
    recognizing.value = true;
    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false }),
        transactionItemsStore.loadAllItems({ force: false }),
        transactionsStore.recognizeReceiptImageByOCR(file)
    ]).then((results) => {
        const ocrResult = results[4] as { transactions: RecognizedReceiptImageResponse[], config?: { hideCategoryColumn?: boolean, hideItemsColumn?: boolean, hideTagsColumn?: boolean } };
        recognizedList.value = ocrResult.transactions;
        // 从 OCR 识别响应中提取配置（仅在首次识别时应用）
        applyOCRConfig(ocrResult.config);
        recognizing.value = false;
    }).catch((error) => {
        recognizing.value = false;
        const msg = error?.message ?? (typeof error === 'string' ? error : 'Recognize failed');
        showToast(msg);
    });
}

function onPreviewRecognize(): void {
    if (!previewFile.value || recognizing.value) return;
    recognizeFile(previewFile.value);
    previewOpened.value = false;
}

function onPreviewSelectImage(): void {
    if (recognizing.value) return;
    triggerFileInput();
}

function onPreviewCancel(): void {
    previewOpened.value = false;
    if (previewImageSrc.value) {
        URL.revokeObjectURL(previewImageSrc.value);
    }
    previewImageSrc.value = null;
    previewFile.value = null;
}

function reset(): void {
    recognizedList.value = [];
    addedRowIndices.value = new Set();
    onPreviewCancel();
}

function onPageAfterIn(): void {
    // 如果是首次进入且没有识别结果，则自动触发选图
    if (!recognizedList.value.length && !recognizing.value && !previewOpened.value) {
        triggerFileInput();
    }

    // 从“添加交易”页返回后，如果有标记的 OCR 行索引，则禁用对应行的“添加”按钮
    const rowIndex = transactionsStore.lastOCRAddedRowIndex;
    if (rowIndex !== null && rowIndex !== undefined) {
        addedRowIndices.value = new Set(addedRowIndices.value).add(rowIndex);
        transactionsStore.setLastOCRAddedRowIndex(null);
    }
}

function buildAddUrl(item: RecognizedReceiptImageResponse, rowIndex?: number): string {
    const params = new URLSearchParams();
    if (item.time != null) params.set('time', String(item.time));
    if (item.type != null) params.set('type', String(item.type));
    if (item.categoryId) params.set('categoryId', item.categoryId);
    if (item.sourceAccountId) params.set('accountId', item.sourceAccountId);
    if (item.destinationAccountId) params.set('destinationAccountId', item.destinationAccountId);
    if (item.sourceAmount != null) params.set('amount', String(item.sourceAmount));
    if (item.destinationAmount != null) params.set('destinationAmount', String(item.destinationAmount));
    if (item.tagIds?.length) params.set('tagIds', item.tagIds.join(','));
    if (item.itemIds?.length) params.set('itemIds', item.itemIds.join(','));
    if (item.comment) params.set('comment', item.comment);
    if (rowIndex !== undefined) params.set('fromOCRRowIndex', String(rowIndex));
    params.set('noTransactionDraft', 'true');
    return `/transaction/add?${params.toString()}`;
}

function onAdd(item: RecognizedReceiptImageResponse, idx: number): void {
    const url = buildAddUrl(item, idx);
    props.f7router.navigate(url);
}
</script>

<style scoped>
.display-none {
    display: none;
}
.ocr-select-block {
    padding-top: 2rem;
    text-align: center;
}
.ocr-hint {
    margin-bottom: 1rem;
    color: var(--f7-block-title-medium-text-color);
}
.ocr-result-item {
    position: relative;
    padding: 3px;
    border-bottom: 1px solid var(--f7-list-strong-border-color);
}
.ocr-result-item:last-child {
    border-bottom: none;
}
.ocr-fields {
    text-align: left;
    padding-right: 50px;
}
.ocr-field {
    display: flex;
    align-items: baseline;
    gap: 0.35rem;
    margin-bottom: 0.25rem;
}
.ocr-label {
    flex-shrink: 0;
    font-size: 0.85em;
    color: var(--f7-block-title-medium-text-color);
}
.ocr-value {
    min-width: 0;
    font-size: 0.9em;
}
.ocr-value.ocr-desc {
    word-break: break-word;
}
.ocr-add-wrap {
    position: absolute;
    right: 16px;
    bottom: 16px;
}
.ocr-add-btn { flex-shrink: 0; }
/* OCR 预览页面样式 - navbar 浮动在图片上方 */
.ocr-preview-page > .page-content {
    padding-top: 0;
    padding-bottom: 0;
    overflow: hidden;
}

/* OCR 预览 navbar 浮动在图片上方，半透明背景 */
.ocr-preview-page .ocr-preview-navbar {
    background: rgba(var(--f7-navbar-bg-color-rgb, 255, 255, 255), 0.85);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ocr-preview-block {
    --f7-block-margin-vertical: 0;
    --f7-block-padding-horizontal: 0;
    --f7-block-padding-vertical: 0;
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 0.5rem;
    padding-bottom: 50px; /* 为底部 navbar 预留空间 */
}
.ocr-preview-img {
    max-width: 100%;
    max-height: 100vh;
    object-fit: contain;
    border-radius: 8px;
}
</style>
