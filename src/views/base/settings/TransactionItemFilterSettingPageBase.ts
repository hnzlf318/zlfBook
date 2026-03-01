import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionItemsStore } from '@/stores/transactionItem.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { entries, keys, values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';
import { DEFAULT_ITEM_GROUP_ID } from '@/consts/item.ts';

import { TransactionItemGroup } from '@/models/transaction_item_group.ts';
import type { TransactionItem } from '@/models/transaction_item.ts';
import { TransactionItemFilter } from '@/models/transaction.ts';

import { objectFieldToArrayItem } from '@/lib/common.ts';

export enum TransactionItemFilterState {
    Default = 0,
    Include = 1,
    Exclude = 2
}

interface TransactionGroupItemFilterTypes {
    includeType: number;
    excludeType: number;
}

function getEmptyGroupItemFilterTypesMap(allTransactionItemsByGroupMap: Record<string, TransactionItem[]>): Record<string, TransactionGroupItemFilterTypes> {
    const ret: Record<string, TransactionGroupItemFilterTypes> = {};

    for (const groupId of keys(allTransactionItemsByGroupMap)) {
        ret[groupId] = {
            includeType: TransactionTagFilterType.HasAny.type,
            excludeType: TransactionTagFilterType.NotHasAny.type
        };
    }

    return ret;
}

export function useTransactionItemFilterSettingPageBase(type?: string) {
    const { tt } = useI18n();

    const transactionItemsStore = useTransactionItemsStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterContent = ref<string>('');

    const itemFilterStateMap = ref<Record<string, TransactionItemFilterState>>({});
    const groupItemFilterTypesMap = ref<Record<string, TransactionGroupItemFilterTypes>>(getEmptyGroupItemFilterTypesMap(transactionItemsStore.allTransactionItemsByGroupMap));

    const lowerCaseFilterContent = computed<string>(() => filterContent.value.toLowerCase());

    const title = computed<string>(() => {
        return tt('Filter Transaction Items');
    });

    const applyText = computed<string>(() => {
        return 'Apply';
    });

    const groupItemFilterStateCountMap = computed<Record<string, Record<TransactionItemFilterState, number>>>(() => {
        const ret: Record<string, Record<TransactionItemFilterState, number>> = {};

        for (const [groupId, items] of entries(transactionItemsStore.allTransactionItemsByGroupMap)) {
            const stateCountMap: Record<TransactionItemFilterState, number> = {
                [TransactionItemFilterState.Default]: 0,
                [TransactionItemFilterState.Include]: 0,
                [TransactionItemFilterState.Exclude]: 0
            };

            for (const item of items) {
                const state = itemFilterStateMap.value[item.id] ?? TransactionItemFilterState.Default;
                stateCountMap[state] = (stateCountMap[state] || 0) + 1;
            }

            ret[groupId] = stateCountMap;
        }

        return ret;
    });

    const allItemGroupsWithDefault = computed<TransactionItemGroup[]>(() => {
        const allGroups: TransactionItemGroup[] = [];
        const itemsInDefaultGroup = transactionItemsStore.allTransactionItemsByGroupMap[DEFAULT_ITEM_GROUP_ID];

        if (itemsInDefaultGroup && itemsInDefaultGroup.length) {
            const defaultGroup = TransactionItemGroup.createNewItemGroup(tt('Default Group'));
            defaultGroup.id = DEFAULT_ITEM_GROUP_ID;
            allGroups.push(defaultGroup);
        }

        for (const itemGroup of transactionItemsStore.allTransactionItemGroups) {
            const itemsInGroup = transactionItemsStore.allTransactionItemsByGroupMap[itemGroup.id];

            if (itemsInGroup && itemsInGroup.length) {
                allGroups.push(itemGroup);
            }
        }

        return allGroups;
    });

    const allVisibleItems = computed<Record<string, TransactionItem[]>>(() => {
        const ret: Record<string, TransactionItem[]> = {};
        const allItemGroups = transactionItemsStore.allTransactionItemsByGroupMap;

        for (const [groupId, items] of entries(allItemGroups)) {
            const visibleItems: TransactionItem[] = [];

            for (const item of items) {
                if (!showHidden.value && item.hidden) {
                    continue;
                }

                if (lowerCaseFilterContent.value && !item.name.toLowerCase().includes(lowerCaseFilterContent.value)) {
                    continue;
                }

                visibleItems.push(item);
            }

            if (visibleItems.length > 0) {
                ret[groupId] = visibleItems;
            }
        }

        return ret;
    });

    const allVisibleItemGroupIds = computed<string[]>(() => objectFieldToArrayItem(allVisibleItems.value));
    const hasAnyAvailableItem = computed<boolean>(() => transactionItemsStore.allAvailableItemsCount > 0);
    const hasAnyVisibleItem = computed<boolean>(() => {
        for (const items of values(allVisibleItems.value)) {
            if (items.length > 0) {
                return true;
            }
        }
        return false;
    });

    function loadFilterItemIds(): boolean {
        let itemFilters: TransactionItemFilter[] = [];

        if (type === 'statisticsCurrent') {
            itemFilters = TransactionItemFilter.parse(statisticsStore.transactionStatisticsFilter.itemFilter);
        } else if (type === 'transactionListCurrent') {
            itemFilters = TransactionItemFilter.parse(transactionsStore.transactionsFilter.itemFilter);
        } else {
            return false;
        }

        const allItemIdsMap: Record<string, TransactionItemFilterState> = {};
        const allGroupItemFilterTypesMap: Record<string, TransactionGroupItemFilterTypes> = getEmptyGroupItemFilterTypesMap(transactionItemsStore.allTransactionItemsByGroupMap);

        for (const transactionItem of values(transactionItemsStore.allTransactionItemsMap)) {
            allItemIdsMap[transactionItem.id] = TransactionItemFilterState.Default;
        }

        for (const itemFilter of itemFilters) {
            let state: TransactionItemFilterState = TransactionItemFilterState.Default;

            if (itemFilter.type === TransactionTagFilterType.HasAny || itemFilter.type === TransactionTagFilterType.HasAll) {
                state = TransactionItemFilterState.Include;
            } else if (itemFilter.type === TransactionTagFilterType.NotHasAny || itemFilter.type === TransactionTagFilterType.NotHasAll) {
                state = TransactionItemFilterState.Exclude;
            } else {
                continue;
            }

            for (const itemId of itemFilter.itemIds) {
                const item = transactionItemsStore.allTransactionItemsMap[itemId];

                if (!item) {
                    continue;
                }

                const groupFilterTypes = allGroupItemFilterTypesMap[item.groupId];

                if (groupFilterTypes) {
                    if (state === TransactionItemFilterState.Include) {
                        groupFilterTypes.includeType = itemFilter.type.type;
                    } else if (state === TransactionItemFilterState.Exclude) {
                        groupFilterTypes.excludeType = itemFilter.type.type;
                    }

                    allItemIdsMap[itemId] = state;
                }
            }
        }

        itemFilterStateMap.value = allItemIdsMap;
        groupItemFilterTypesMap.value = allGroupItemFilterTypesMap;
        return true;
    }

    function saveFilterItemIds(): boolean {
        const itemFilters: TransactionItemFilter[] = [];
        let changed = true;

        for (const [groupId, items] of entries(transactionItemsStore.allTransactionItemsByGroupMap)) {
            const groupFilterTypes = groupItemFilterTypesMap.value[groupId];

            if (groupFilterTypes && items && items.length > 0) {
                const includeItemIds: string[] = [];
                const excludeItemIds: string[] = [];

                for (const item of items) {
                    const state = itemFilterStateMap.value[item.id] ?? TransactionItemFilterState.Default;

                    if (state === TransactionItemFilterState.Include) {
                        includeItemIds.push(item.id);
                    } else if (state === TransactionItemFilterState.Exclude) {
                        excludeItemIds.push(item.id);
                    }
                }

                if (includeItemIds.length > 0) {
                    const includeItemFilter = TransactionItemFilter.create(includeItemIds, TransactionTagFilterType.parse(groupFilterTypes.includeType) ?? TransactionTagFilterType.HasAny);
                    itemFilters.push(includeItemFilter);
                }

                if (excludeItemIds.length > 0) {
                    const excludeItemFilter = TransactionItemFilter.create(excludeItemIds, TransactionTagFilterType.parse(groupFilterTypes.excludeType) ?? TransactionTagFilterType.NotHasAny);
                    itemFilters.push(excludeItemFilter);
                }
            }
        }

        if (type === 'statisticsCurrent') {
            changed = statisticsStore.updateTransactionStatisticsFilter({
                itemFilter: TransactionItemFilter.toTextualItemFilters(itemFilters)
            });

            if (changed) {
                statisticsStore.updateTransactionStatisticsInvalidState(true);
            }
        } else if (type === 'transactionListCurrent') {
            changed = transactionsStore.updateTransactionListFilter({
                itemFilter: TransactionItemFilter.toTextualItemFilters(itemFilters)
            });

            if (changed) {
                transactionsStore.updateTransactionListInvalidState(true);
            }
        }

        return changed;
    }

    return {
        // states
        loading,
        showHidden,
        filterContent,
        itemFilterStateMap,
        groupItemFilterTypesMap,
        // computed states
        title,
        applyText,
        groupItemFilterStateCountMap,
        allItemGroupsWithDefault,
        allVisibleItems,
        allVisibleItemGroupIds,
        hasAnyAvailableItem,
        hasAnyVisibleItem,
        // functions
        loadFilterItemIds,
        saveFilterItemIds
    };
}

