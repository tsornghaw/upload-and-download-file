import React, { ChangeEvent, ChangeEventHandler, useEffect, useRef, useState } from 'react';

type DataItem = {
    dataname: string;
    download_url: string;
    share_limti: number;
  };

const Home = (props: { name: string }) => {
    
    // Define a state variable to store the options fetched from the API
    const [fileList, setFileList] = useState<FileList | null>(null);
    const [selectedImage, setSelectedImage] = useState('');
    const [options, setOptions] = useState<DataItem[]>([]);
    const [selectedOption, setselectedOption] = useState<string>("");
    const [downloadmessage, setdownloadMessage] = useState('');
    const [uploadmessage, setuploadMessage] = useState('');
    const [downloadTimes, setDownloadTimes] = useState(5);

    const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
        setFileList(e.target.files);
        
        const selectedFiles = e.target.files

        if (selectedFiles && selectedFiles.length > 0) {
            const file = selectedFiles[0];
            setSelectedImage(URL.createObjectURL(file));
        }
    };

    const handleDownload = async () => {
        if (selectedOption) {
            const download_url = 'http://localhost:8000/api/auth/download/' + selectedOption
            try {
                const response = await fetch(download_url, {
                    credentials: 'include',
                });

                if (response.ok) {
                    // The download was successful, set a success message
                    setdownloadMessage('File download successful');
                } else {
                    // The download failed, set an error message
                    setdownloadMessage('Error downloading the file: ' + response.statusText);
                }
                console.log(response)
            } catch (error) {
                console.error('Error downloading the file:', error);
            }
        } else {
            console.error('No option selected for download');
        }
    };

    const handleUpload = () => {
        if (fileList === null || fileList.length === 0) {
            return
        }
        const formData = new FormData()
        formData.append(`file`, fileList[0], fileList[0].name)
        formData.append('downloadTimes', downloadTimes.toString());
        console.log(downloadTimes)

        fetch('http://localhost:8000/api/auth/upload', {
            method: 'POST',
            body: formData,
            credentials: 'include',
        })
        .then((response) => {
            if (response.ok) {
                // Upload successful, set success message
                setuploadMessage('File uploaded successfully');
                return response.json();
            } else {
                throw new Error('Upload failed');
            }
          })

        // Re-fetch data and update options when the button is clicked
        fetch('http://localhost:8000/api/auth/UserSearchAllData', {
            credentials: 'include',
        })
        .then((response) => response.json())
        .then((data) => {
            setOptions(data);
        })
        .catch((error) => {
            console.error('Error fetching data from the API:', error);
        });
    }

    // Use the useEffect hook to fetch data from the API when the component mounts
    useEffect(() => {
        fetch('http://localhost:8000/api/auth/UserSearchAllData', {
            credentials: 'include',
        })
        .then((response) => response.json())
        .then((data) => {
            setOptions(data);
        })
        .catch((error) => {
            console.error('Error fetching data from the API:', error);
        });
    }, []); // The empty dependency array ensures this effect runs only once on component mount

    // Handle the selection change
    const handleSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setselectedOption(event.target.value);
    };

    const handleRangeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        // Update the state with the current range value
        setDownloadTimes(parseInt(event.target.value, 10));
    };

    
    return (
        <div>
            {props.name ? 
                (
                    <React.Fragment>
                        {/* If Authenticated - Upload */}
                        <input type="file" id="myFile" name="filename"
                            accept="image/*,application/pdf,.doc,.ppt,.csv,.png,.jpg,.pdf,.xls,application/msword,application/vnd.ms-powerpoint"
                            onChange={handleFileChange}/>
                        {selectedImage && (
                            <img
                            src={selectedImage}
                            alt="Selected"
                            style={{ maxWidth: '600px', maxHeight: '600px' }} // Set your desired max width and max height
                            />
                        )}
                        <br />
                        <label htmlFor="vol">Download Times (between 0 and 10):</label>
                        <input type="range" id="vol" name="vol" min="0" max="10" value={downloadTimes} onChange={handleRangeChange}></input>
                        <p>Selected Value: {downloadTimes}</p>
                        <br />
                        <button onClick={handleUpload}>Submit</button>
                        <p>{uploadmessage}</p>
                        <br />

                        {/* If Authenticated - Download */}
                        <label>Select an option: </label>
                        <select value={selectedOption} onChange={handleSelectChange}>
                        <option value="">Select an option</option>
                        {options && options.map((option, index) => (
                            <option key={option.dataname + index} value={option.download_url}>
                            {option.dataname}
                            </option>
                        ))}
                        </select>
                        {selectedOption && <p>Selected option: {selectedOption}</p>}
                        <button onClick={handleDownload}>Download File</button>
                        <p>{downloadmessage}</p>
                        <br  />
                    </React.Fragment>
                )
                :
                (
                    // If not Authenticate
                    'You are not logged in'
                )
            }
        </div>
    );
}

export default Home;


function setSelectedImage(arg0: string) {
    throw new Error('Function not implemented.');
}
// Optimization: pjchender - Day23