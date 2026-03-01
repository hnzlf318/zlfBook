<template>
    <v-card :class="{ 'pa-sm-1 pa-md-2': dialogMode }">
        <template #title>
            <v-row>
                <v-col cols="6">
                    <div :class="{ 'text-h4': dialogMode, 'text-wrap': true }">
                        {{ tt(title) }}
                    </div>
                </v-col>
                <v-col cols="6" class="d-flex align-center">
                    <v-spacer v-if="!dialogMode"/>
                    <v-text-field density="compact" :disabled="loading || !hasAnyAvailableItem"
                                  :prepend-inner-icon="mdiMagnify"
                                  :placeholder="tt('Find item')"
                                  v-model="filterContent"
                                  v-if="dialogMode"></v-text-field>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :disabled="loading || !hasAnyAvailableItem" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Set All to Included')"
                                             :disabled="!hasAnyVisibleItem"
                                             @click="setAllItemsState(TransactionItemFilterState.Include)"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Set All to Default')"
                                             :disabled="!hasAnyVisibleItem"
                                             @click="setAllItemsState(TransactionItemFilterState.Default)"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Set All to Excluded')"
                                             :disabled="!hasAnyVisibleItem"
                                             @click="setAllItemsState(TransactionItemFilterState.Exclude)"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiEyeOutline"
                                             :title="tt('Show Hidden Transaction Items')"
                                             v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                <v-list-item :prepend-icon="mdiEyeOffOutline"
                                             :title="tt('Hide Hidden Transaction Items')"
                                             v-if="showHidden" @click="showHidden = false"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </v-col>
            </v-row>
        </template>

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text v-if="!loading && !hasAnyVisibleItem">
            <span class="text-body-1">{{ tt('No available item') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'flex-grow-1 overflow-y-auto': dialogMode }" v-else-if="!loading && hasAnyVisibleItem">
            <v-expansion-panels class="item-categories" multiple v-model="expandItemGroups">
                <template :key="itemGroup.id" v-for="itemGroup in allItemGroupsWithDefault">
                    <v-expansion-panel class="border" :value="itemGroup.id" v-if="allVisibleItems[itemGroup.id] && allVisibleItems[itemGroup.id]!.length > 0">
                        <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                            <span class="ms-3 text-truncate">{{ itemGroup.name }}</span>
                            <v-spacer/>
                            <div class="d-flex me-3" v-if="groupItemFilterTypesMap[itemGroup.id] && groupItemFilterStateCountMap[itemGroup.id]">
                                <v-btn color="secondary" density="compact" variant="outlined"
                                       v-if="groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Include] && groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Include] > 1">
                                    {{ groupItemFilterTypesMap[itemGroup.id]!.includeType === TransactionTagFilterType.HasAll.type ? tt(TransactionTagFilterType.HasAll.name) : tt(TransactionTagFilterType.HasAny.name) }}
                                    <v-menu activator="parent">
                                        <v-list>
                                            <v-list-item :key="filterType.type" :title="tt(filterType.name)"
                                                         :append-icon="groupItemFilterTypesMap[itemGroup.id]!.includeType === filterType.type ? mdiCheck : undefined"
                                                         v-for="filterType in [TransactionTagFilterType.HasAny, TransactionTagFilterType.HasAll]"
                                                         @click="updateTransactionItemGroupIncludeType(itemGroup, filterType)"></v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                                       v-if="groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Exclude] && groupItemFilterStateCountMap[itemGroup.id]![TransactionItemFilterState.Exclude] > 1">
                                    {{ groupItemFilterTypesMap[itemGroup.id]!.excludeType === TransactionTagFilterType.NotHasAll.type ? tt(TransactionTagFilterType.NotHasAll.name) : tt(TransactionTagFilterType.NotHasAny.name) }}
                                    <v-menu activator="parent">
                                        <v-list>
                                            <v-list-item :key="filterType.type" :title="tt(filterType.name)"
                                                         :append-icon="groupItemFilterTypesMap[itemGroup.id]!.excludeType === filterType.type ? mdiCheck : undefined"
                                                         v-for="filterType in [TransactionTagFilterType.NotHasAny, TransactionTagFilterType.NotHasAll]"
                                                         @click="updateTransactionItemGroupExcludeType(itemGroup, filterType)"></v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                            </div>
                        </v-expansion-panel-title>
                        <v-expansion-panel-text>
                            <v-list rounded density="comfortable" class="pa-0">
                                <template :key="transactionItem.id"
                                          v-for="transactionItem in allVisibleItems[itemGroup.id]">
                                    <v-list-item class="ps-2">
                                        <template #prepend>
                                            <v-badge class="right-bottom-icon" color="secondary"
                                                     location="bottom right" offset-x="2" offset-y="2" :icon="mdiEyeOffOutline"
                                                     v-if="transactionItem.hidden">
                                                <v-icon size="24" :icon="mdiFormatListBulleted"/>
                                            </v-badge>
                                            <v-icon size="24" :icon="mdiFormatListBulleted" v-else-if="!transactionItem.hidden"/>
                                            <span class="ms-3">{{ transactionItem.name }}</span>
                                        </template>
                                        <template #append>
                                            <v-btn-toggle class="toggle-buttons" density="compact" variant="outlined"
                                                          mandatory="force" divided
                                                          :model-value="itemFilterStateMap[transactionItem.id]"
                                                          @update:model-value="updateTransactionItemState(transactionItem, $event)">
                                                <v-btn :value="TransactionItemFilterState.Include">{{ tt('Included') }}</v-btn>
                                                <v-btn :value="TransactionItemFilterState.Default">{{ tt('Default') }}</v-btn>
                                                <v-btn :value="TransactionItemFilterState.Exclude">{{ tt('Excluded') }}</v-btn>
                                            </v-btn-toggle>
                                        </template>
                                    </v-list-item>
                                </template>
                            </v-list>
                        </v-expansion-panel-text>
                    </v-expansion-panel>
                </template>
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                <v-btn :disabled="!hasAnyAvailableItem" @click="save">{{ tt(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    useTransactionItemFilterSettingPageBase,
    TransactionItemFilterState
} from '@/views/base/settings/TransactionItemFilterSettingPageBase.ts';

import { useTransactionItemsStore } from '@/stores/transactionItem.ts';

import { values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';

import type { TransactionItemGroup } from '@/models/transaction_item_group.ts';
import type { TransactionItem } from '@/models/transaction_item.ts';

import {
    mdiMagnify,
    mdiCheck,
    mdiSelectAll,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical,
    mdiFormatListBulleted
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    type: string;
    dialogMode?: boolean;
    autoSave?: boolean;
}>();

const emit = defineEmits<{
    (e: 'settings:change', changed: boolean): void;
}>();

const { tt } = useI18n();

const {
    loading,
    showHidden,
    filterContent,
    itemFilterStateMap,
    groupItemFilterTypesMap,
    title,
    applyText,
    groupItemFilterStateCountMap,
    allItemGroupsWithDefault,
    allVisibleItems,
    allVisibleItemGroupIds,
    hasAnyAvailableItem,
    hasAnyVisibleItem,
    loadFilterItemIds,
    saveFilterItemIds
} = useTransactionItemFilterSettingPageBase(props.type);

const transactionItemsStore = useTransactionItemsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandItemGroups = ref<string[]>(allVisibleItemGroupIds.value);

function init(): void {
    transactionItemsStore.loadAllItems({
        force: false
    }).then(() => {
        loading.value = false;
        expandItemGroups.value = allVisibleItemGroupIds.value;

        if (!loadFilterItemIds()) {
            snackbar.value?.showError('Parameter Invalid');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function updateTransactionItemState(transactionItem: TransactionItem, value: TransactionItemFilterState): void {
    itemFilterStateMap.value[transactionItem.id] = value;

    if (props.autoSave) {
        save();
    }
}

function updateTransactionItemGroupIncludeType(itemGroup: TransactionItemGroup, filterType: TransactionTagFilterType): void {
    const itemFilterTypes = groupItemFilterTypesMap.value[itemGroup.id];

    if (!itemFilterTypes) {
        return;
    }

    itemFilterTypes.includeType = filterType.type;

    if (props.autoSave) {
        save();
    }
}

function updateTransactionItemGroupExcludeType(itemGroup: TransactionItemGroup, filterType: TransactionTagFilterType): void {
    const itemFilterTypes = groupItemFilterTypesMap.value[itemGroup.id];

    if (!itemFilterTypes) {
        return;
    }

    itemFilterTypes.excludeType = filterType.type;

    if (props.autoSave) {
        save();
    }
}

function setAllItemsState(value: TransactionItemFilterState): void {
    for (const items of values(allVisibleItems.value)) {
        for (const item of items) {
            itemFilterStateMap.value[item.id] = value;
        }
    }

    if (props.autoSave) {
        save();
    }
}

function save(): void {
    const changed = saveFilterItemIds();
    emit('settings:change', changed);
}

function cancel(): void {
    emit('settings:change', false);
}

init();
</script>

<style>
.item-categories .item-filter-state-toggle {
    overflow-x: auto;
    white-space: nowrap;
}

.item-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 0;
    padding-inline-start: 20px;
}

.item-categories .v-expansion-panel--active:not(:first-child),
.item-categories .v-expansion-panel--active + .v-expansion-panel {
    margin-top: 1rem;
}
</style>

