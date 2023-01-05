import React from 'react';
import { FileStore } from './FileStore/FileStore';
import { FileStoreHTTP } from './FileStore/FileStoreHTTP';
import usePromise from './hooks/usePromise';
import useStorageState from './hooks/useStorageState';
import TextInput from './TextInput';
import UploadFilesInput from './UploadFilesInput';
import DownloadFileButton from './DownloadFileButton';
import DownloadAllFilesButton from './DownloadAllFilesButton';

function App() {
  // TODO: this appears to be getting called twice on initial load
  const [baseURL, setBaseURL] = useStorageState(localStorage, "baseURL", "http://localhost:9000")
  const fileStore: FileStore = React.useMemo(() => new FileStoreHTTP(baseURL), [baseURL])
  const [lastUpload, setLastUpload] = React.useState(0)
  const onUpdate = React.useCallback(() => {
    setLastUpload(Date.now())
  }, [])
  const listFilesPromise = React.useMemo(() => fileStore.listFiles(), [fileStore, lastUpload])
  const listFilesState = usePromise(listFilesPromise)

  return (
    <div className="container">
      <form>
        <div className="row">
          <div className="col-md">
            <div className="form-group">
              <label htmlFor="input-baseURL">Base URL</label>
              <TextInput value={baseURL} onChange={setBaseURL} id="input-baseURL" placeholder="http://localhost:9000" className="form-control" />
            </div>
          </div>
          <div className="col-md">
            <div>
              <label htmlFor="input-files">Upload Files</label>
              <UploadFilesInput
                fileStore={fileStore}
                onUpdate={onUpdate}
                id="input-files"
                className="form-control"
              />
            </div>
          </div>
        </div>
      </form>
      {listFilesState.type === 'PENDING' && <p>Loading...</p>}
      {listFilesState.type === 'LOADED' && (
        <React.Fragment>
          <table className="table">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th scope="col">Actions</th>
              </tr>
            </thead>
            <tbody>
              {listFilesState.data.map((fileInfo) => (
                <tr key={fileInfo.id}>
                  <td>{fileInfo.name}</td>
                  <td>
                    <DownloadFileButton fileStore={fileStore} fileInfo={fileInfo} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <DownloadAllFilesButton fileStore={fileStore} fileInfos={listFilesState.data} />
        </React.Fragment>
      )}
      {listFilesState.type === 'ERROR' && <p>{listFilesState.error.message}</p>}
    </div>
  );
}

export default App;
