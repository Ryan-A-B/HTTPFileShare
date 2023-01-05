export interface FileInfo {
    id: string
    name: string
    mime_type: string
}

export interface FileStore {
    addFile(file: File): Promise<FileInfo>
    listFiles(): Promise<FileInfo[]>
    downloadFile(id: string): Promise<Blob>
}
