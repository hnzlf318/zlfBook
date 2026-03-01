<template>
    <f7-page class="page-navbar-bottom" with-subnavbar @page:beforein="onPageBeforeIn" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right :class="{ 'navbar-compact-icons': true, 'disabled': loading }">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableItem }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !hasAnyAvailableItem }" @click="save"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :inner="false">
                <f7-searchbar
                    custom-searchs
                    :class="{ 'disabled': loading }"
                    :value="filterContent"
                    :placeholder="tt('Find item')"
                    :disable-button-text="tt('Cancel')"
                    @change="filterContent = $event.target.value"
                ></f7-searchbar>
            </f7-subnavbar>
        </f7-navbar>

        <f7-block class="combination-list-wrapper margin-vertical skeleton-text" v-if="loading">
            <f7-accordion-item>
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers
                                 class="combination-list-header combination-list-opened">
                            <f7-list-item group-title>
                                <small>{{ tt('Transaction Items') }}</small>
                                <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content style="height: auto">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item checkbox class="disabled" :title="tt('Transaction Item Name')"
                                      :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                            <template #media>
                                <f7-icon f7="app_fill"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-list strong inset dividers accordion-list class="margin-top" v-if="!loading && !hasAnyVisibleItem">
            <f7-list-item :title="tt('No available item')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="itemGroup.id" v-for="itemGroup in allItemGroupsWithDefault"
                  v-show="!loading && hasAnyVisibleItem">
            <f7-accordion-item :opened="collapseStates[itemGroup.id]?.opened ?? true"
                               @accordion:open="collapseStates[itemGroup.id]!.opened = true"
                               @accordion:close="collapseStates[itemGroup.id]!.opened = false"
                               v-if="allVisibleItems[itemGroup.id] && allVisibleItems[itemGroup.id]!.length > 0">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers
                                 class="combination-list-header"
                                 :class="collapseStates[itemGroup.id]?.opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item group-title>
                                <small class="item-group-title">{{ itemGroup.name }}</small>
                                <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[itemGroup.id]?.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[itemGroup.id]?.opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item link="#"
                                      popover-open=".item-filter-include-type-popover-menu"
                                      :title="tt(TransactionTagFilterType.parse(groupItemFilterTypesMap[itemGroup.id]?.includeType as number)?.name as string)"
                                      @click="currentTransactionItemGroupId = itemGroup.id"
                                      v-if="groupItemFilterTypesMap[itemGroup.id] && groupItemFilterStateCountMap[itemGroup.id] && groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Include] && groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Include] > 1 && TransactionTagFilterType.parse(groupItemFilterTypesMap[itemGroup.id]?.includeType as number)">
                        </f7-list-item>
                        <f7-list-item link="#"
                                      popover-open=".item-filter-exclude-type-popover-menu"
                                      :title="tt(TransactionTagFilterType.parse(groupItemFilterTypesMap[itemGroup.id]?.excludeType as number)?.name as string)"
                                      @click="currentTransactionItemGroupId = itemGroup.id"
                                      v-if="groupItemFilterTypesMap[itemGroup.id] && groupItemFilterStateCountMap[itemGroup.id] && groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Exclude] && groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Exclude] > 1 && TransactionTagFilterType.parse(groupItemFilterTypesMap[itemGroup.id]?.excludeType as number)">
                        </f7-list-item>
                        <f7-list-item link="#"
                                      popover-open=".item-filter-state-popover-menu"
                                      :title="transactionItem.name"
                                      :value="transactionItem.id"
                                      :key="transactionItem.id"
                                      :after="tt(itemFilterStateMap[transactionItem.id] === TransactionItemFilterState.Include ? 'Included' : itemFilterStateMap[transactionItem.id] === TransactionItemFilterState.Exclude ? 'Excluded' : 'Default')"
                                      v-for="transactionItem in allVisibleItems[itemGroup.id]"
                                      v-show="showHidden || !transactionItem.hidden"
                                      @click="currentTransactionItemId = transactionItem.id">
                            <template #media>
                                <f7-icon class="transaction-item-icon" f7="list_bullet">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="transactionItem.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-popover class="item-filter-include-type-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt(filterType.name)"
                              :class="{ 'list-item-selected': groupItemFilterTypesMap[currentTransactionItemGroupId]?.includeType === filterType.type }"
                              :key="filterType.type"
                              v-for="filterType in [TransactionTagFilterType.HasAny, TransactionTagFilterType.HasAll]"
                              @click="updateTransactionItemGroupIncludeType(filterType)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="groupItemFilterTypesMap[currentTransactionItemGroupId]?.includeType === filterType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="item-filter-exclude-type-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt(filterType.name)"
                              :class="{ 'list-item-selected': groupItemFilterTypesMap[currentTransactionItemGroupId]?.excludeType === filterType.type }"
                              :key="filterType.type"
                              v-for="filterType in [TransactionTagFilterType.NotHasAny, TransactionTagFilterType.NotHasAll]"
                              @click="updateTransactionItemGroupExcludeType(filterType)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="groupItemFilterTypesMap[currentTransactionItemGroupId]?.excludeType === filterType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="item-filter-state-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="state.displayName"
                              :class="{ 'list-item-selected': itemFilterStateMap[currentTransactionItemId] === state.type }"
                              :key="state.type"
                              v-for="state in [
                                  { type: TransactionItemFilterState.Include, displayName: tt('Included') },
                                  { type: TransactionItemFilterState.Default, displayName: tt('Default') },
                                  { type: TransactionItemFilterState.Exclude, displayName: tt('Excluded') }
                              ]"
                              @click="updateCurrentTransactionItemState(state.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="itemFilterStateMap[currentTransactionItemId] === state.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleItem }" @click="setAllItemsState(TransactionItemFilterState.Include)">{{ tt('Set All to Included') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleItem }" @click="setAllItemsState(TransactionItemFilterState.Default)">{{ tt('Set All to Default') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleItem }" @click="setAllItemsState(TransactionItemFilterState.Exclude)">{{ tt('Set All to Excluded') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Transaction Items') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Transaction Items') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import {
    useTransactionItemFilterSettingPageBase,
    TransactionItemFilterState
} from '@/views/base/settings/TransactionItemFilterSettingPageBase.ts';

import { useTransactionItemsStore } from '@/stores/transactionItem.ts';

import { values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';

interface CollapseState {
    opened: boolean;
}

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const query = props.f7route.query;

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const {
    loading,
    showHidden,
    filterContent,
    itemFilterStateMap,
    groupItemFilterTypesMap,
    title,
    groupItemFilterStateCountMap,
    allItemGroupsWithDefault,
    allVisibleItems,
    allVisibleItemGroupIds,
    hasAnyAvailableItem,
    hasAnyVisibleItem,
    loadFilterItemIds,
    saveFilterItemIds
} = useTransactionItemFilterSettingPageBase(query['type']);

const transactionItemsStore = useTransactionItemsStore();

const loadingError = ref<unknown | null>(null);
const currentTransactionItemGroupId = ref<string>('');
const currentTransactionItemId = ref<string>('');
const showMoreActionSheet = ref<boolean>(false);

const collapseStates = ref<Record<string, CollapseState>>(getInitCollapseState(allVisibleItemGroupIds.value));

function getInitCollapseState(itemGroupIds: string[]): Record<string, CollapseState> {
    const states: Record<string, CollapseState> = {};

    for (const itemGroupId of itemGroupIds) {
        states[itemGroupId] = {
            opened: true
        };
    }

    return states;
}

function init(): void {
    transactionItemsStore.loadAllItems({
        force: false
    }).then(() => {
        loading.value = false;
        collapseStates.value = getInitCollapseState(allVisibleItemGroupIds.value);

        if (!loadFilterItemIds()) {
            showToast('Parameter Invalid');
            loadingError.value = 'Parameter Invalid';
        }
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function updateTransactionItemGroupIncludeType(filterType: TransactionTagFilterType): void {
    const itemFilterTypes = groupItemFilterTypesMap.value[currentTransactionItemGroupId.value];

    if (!itemFilterTypes) {
        return;
    }

    itemFilterTypes.includeType = filterType.type;
}

function updateTransactionItemGroupExcludeType(filterType: TransactionTagFilterType): void {
    const itemFilterTypes = groupItemFilterTypesMap.value[currentTransactionItemGroupId.value];

    if (!itemFilterTypes) {
        return;
    }

    itemFilterTypes.excludeType = filterType.type;
}

function updateCurrentTransactionItemState(state: number): void {
    itemFilterStateMap.value[currentTransactionItemId.value] = state;
    currentTransactionItemId.value = '';
}

function setAllItemsState(value: TransactionItemFilterState): void {
    for (const items of values(allVisibleItems.value)) {
        for (const item of items) {
            itemFilterStateMap.value[item.id] = value;
        }
    }
}

function save(): void {
    saveFilterItemIds();
    props.f7router.back();
}

function onPageBeforeIn(): void {
    filterContent.value = '';
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.item-group-title {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>

