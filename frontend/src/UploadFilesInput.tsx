import React from 'react'
import { FileStore } from './FileStore/FileStore'

interface Props {
    fileStore: FileStore
    onUpdate: () => void
    id: string
    className: string
}

const UploadFilesInput: React.FunctionComponent<Props> = ({ fileStore, onUpdate, ...props }) => {
    const onChange = React.useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        const files = event.target.files
        if (files === null) return
        for (let i = 0; i < files.length; i++) {
            const file = files.item(i)
            if (file === null) throw new Error("file should not be null")
            fileStore.addFile(file)
        }
        onUpdate()
    }, [fileStore, onUpdate])
    return (
        <input type="file" value="" onChange={onChange} multiple {...props} />
    )
}

export default UploadFilesInput