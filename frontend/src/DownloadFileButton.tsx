import React from 'react'
import { FileStore, FileInfo } from './FileStore/FileStore'

interface Props {
  fileStore: FileStore
  fileInfo: FileInfo
}

const DownloadFileButton: React.FunctionComponent<Props> = ({ fileStore, fileInfo }) => {
  const onClick = React.useCallback(async () => {
    const blob = await fileStore.downloadFile(fileInfo.id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileInfo.name
    a.click()
  }, [fileStore, fileInfo])
  return (
    <button type="button" className="btn btn-primary" onClick={onClick}>
      Download
    </button>
  )
}

export default DownloadFileButton
