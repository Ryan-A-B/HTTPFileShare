import React from 'react'
import { FileStore, FileInfo } from './FileStore/FileStore'

interface Props {
    fileStore: FileStore
    fileInfos: FileInfo[]
}

const DownloadAllFilesButton: React.FunctionComponent<Props> = ({ fileStore, fileInfos }) => {
    const onClick = React.useCallback(() => {
        fileInfos.forEach(async (fileInfo) => {
            const blob = await fileStore.downloadFile(fileInfo.id)
            const url = URL.createObjectURL(blob)
            const a = document.createElement('a')
            a.href = url
            a.download = fileInfo.name
            a.click()
        })
    }, [fileStore, fileInfos])
    return (
        <button type="button" className="btn btn-primary" onClick={onClick}>
            Download All
        </button>
    )
}

export default DownloadAllFilesButton
