<template>
    <v-autocomplete
        item-title="name"
        item-value="id"
        auto-select-first
        persistent-placeholder
        multiple
        chips
        :density="density"
        :variant="variant"
        :closable-chips="!readonly"
        :readonly="readonly"
        :disabled="disabled"
        :label="showLabel ? tt('Transaction Items') : undefined"
        :placeholder="tt('None')"
        :items="allItemsWithGroupHeader"
        :model-value="modelValue"
        v-model:search="itemSearchContent"
        @update:modelValue="updateModelValue"
    >
        <template #chip="{ props, item }">
            <v-chip :prepend-icon="mdiFormatListBulleted" :text="item.title" v-bind="props"/>
        </template>

        <template #subheader="{ props }">
            <v-list-subheader>{{ props['title'] }}</v-list-subheader>
        </template>

        <template #item="{ props, item }">
            <v-list-subheader v-if="item.raw && typeof item.raw === 'object' && 'type' in item.raw && item.raw.type === 'subheader'">
                {{ (item.raw as { title: string }).title }}
            </v-list-subheader>
            <v-list-item :value="item.value" v-bind="props" v-else-if="item.raw instanceof TransactionItem && !item.raw.hidden">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <v-icon size="20" start :icon="mdiFormatListBulleted"/>
                            <span>{{ item.title }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
            <v-list-item :disabled="true" v-bind="props" v-else-if="item.raw instanceof TransactionItem && item.raw.hidden">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <v-icon size="20" start :icon="mdiFormatListBulleted"/>
                            <span>{{ item.title }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
        </template>

        <template #no-data>
            <v-list class="py-0">
                <v-list-item>{{ tt('No available item') }}</v-list-item>
            </v-list>
        </template>
    </v-autocomplete>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { type CommonTransactionItemSelectionProps, useTransactionItemSelectionBase } from '@/components/base/TransactionItemSelectionBase.ts';

import { TransactionItem } from '@/models/transaction_item.ts';

import type { ComponentDensity, InputVariant } from '@/lib/ui/desktop.ts';

import { mdiFormatListBulleted } from '@mdi/js';

interface DesktopTransactionItemSelectionProps extends CommonTransactionItemSelectionProps {
    density?: ComponentDensity;
    variant?: InputVariant;
    readonly?: boolean;
    disabled?: boolean;
    showLabel?: boolean;
}

const props = defineProps<DesktopTransactionItemSelectionProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void;
}>();

const { tt } = useI18n();

const {
    itemSearchContent,
    allItemsWithGroupHeader
} = useTransactionItemSelectionBase(props);

function updateModelValue(newValue: string[]) {
    emit('update:modelValue', newValue);
}
</script>
