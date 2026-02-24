export interface RecognizedReceiptImageResponse {
    readonly type: number;
    readonly time?: number;
    readonly categoryId?: string;
    readonly sourceAccountId?: string;
    readonly destinationAccountId?: string;
    readonly sourceAmount?: number;
    readonly destinationAmount?: number;
    readonly tagIds?: string[];
    readonly itemIds?: string[];
    readonly comment?: string;
}

export interface OCRBillRecognitionConfig {
    readonly hideCategoryColumn: boolean;
    readonly hideItemsColumn: boolean;
    readonly hideTagsColumn: boolean;
    readonly dialogMaxWidth: number;
}

export interface RecognizedReceiptImageListResponse {
    readonly transactions: RecognizedReceiptImageResponse[];
    readonly config?: OCRBillRecognitionConfig;
}
