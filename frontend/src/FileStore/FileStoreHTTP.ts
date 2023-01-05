import { FileInfo, FileStore } from './FileStore'

export class FileStoreHTTP implements FileStore {
    private readonly url: string

    constructor(url: string) {
        this.url = url
    }

    async addFile(file: File): Promise<FileInfo> {
        const formData = new FormData()
        formData.append("file", file, file.name)
        const response = await fetch(this.url + "/files", {
            method: "POST",
            body: formData,
        })
        if (!response.ok) {
            throw new Error("failed to upload file")
        }
        const json = await response.json()
        return json
    }

    async listFiles(): Promise<FileInfo[]> {
        const response = await fetch(this.url + "/files")
        if (!response.ok) {
            throw new Error("failed to list files")
        }
        const json = await response.json()
        return json["items"]
    }

    async downloadFile(id: string): Promise<Blob> {
        const response = await fetch(this.url + "/files/" + id)
        if (!response.ok) {
            throw new Error("failed to download file")
        }
        return response.blob()
    }
}