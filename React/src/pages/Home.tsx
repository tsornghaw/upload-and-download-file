import React, { ChangeEvent, ChangeEventHandler, useEffect, useState } from 'react';
import axios from 'axios';

// Define a type for your data
type DataItem = {
    dataname: string;
    download_url: string;
  };

const Home = (props: { name: string }) => {
    const [fileList, setFileList] = useState<FileList | null>(null);
    const [selectedImage, setSelectedImage] = useState('');

    // Define a state variable to store the options fetched from the API
    const [options, setOptions] = useState<DataItem[]>([]);

    // Define a state variable to store the selected option
    const [selectedOption, setSelectedOption] = useState<string>("");

    const [downloadmessage, setdownloadMessage] = useState('');
    const [uploadmessage, setuploadMessage] = useState('');

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
            // Implement the API call to download here
            console.log('Downloading : ' + selectedOption);
            const download_url = 'http://localhost:8000/api/auth/download/' + selectedOption
            // You can use fetch or any other library to make the API call
            try {
                const response = await fetch(download_url, {
                    credentials: 'include',
                });
                //const blob = await response.blob();

                // Create a URL for the blob data and trigger a download
                // const url = window.URL.createObjectURL(blob);
                // const a = document.createElement('a');
                // a.href = url;
                // a.download = 'test.png'; // Set the desired filename
                // document.body.appendChild(a);
                // a.click();
                // window.URL.revokeObjectURL(url);

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
        console.log(fileList[0])
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
            // Assuming the API response is an array of objects with 'value' and 'label' properties
            setOptions(data);
            console.log(data)
        })
        .catch((error) => {
            console.error('Error fetching data from the API:', error);
        });
    }

    // Use the useEffect hook to fetch data from the API when the component mounts
    useEffect(() => {
        // Replace 'your-api-endpoint' with the actual API endpoint URL
        fetch('http://localhost:8000/api/auth/UserSearchAllData', {
            credentials: 'include',
        })
        .then((response) => response.json())
        .then((data) => {
            // Assuming the API response is an array of objects with 'value' and 'label' properties
            setOptions(data);
            console.log(data)
        })
        .catch((error) => {
            console.error('Error fetching data from the API:', error);
        });
    }, []); // The empty dependency array ensures this effect runs only once on component mount

    // Handle the selection change
    const handleSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setSelectedOption(event.target.value);
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
                        {/* <input type="submit" /> */}
                        <button onClick={handleUpload}>Submit</button>
                        <p>{uploadmessage}</p>
                        <br />


                        {/* If Authenticated - Download */}
                        <label>Select an option: </label>
                        <select value={selectedOption} onChange={handleSelectChange}>
                            {/* <option value="">Select an option</option>
                            {options.map((option) => (
                            <option key={option.dataname} value={option.download_url}>
                                {option.dataname}
                            </option>
                            ))} */}

                        <option value="">Select an option</option>
                        {options.map((option, index) => (
                            <option key={option.dataname + index} value={option.download_url}>
                            {option.dataname}
                            </option>
                        ))}
                        </select>
                        {selectedOption && <p>Selected option: {selectedOption}</p>}
                        <button onClick={handleDownload}>Download File</button>
                        <p>{downloadmessage}</p>

                        <br  />
                        {/* <input
                            type="text"
                            placeholder="Enter file name"
                            value={fileName}
                            onChange={(e) => setFileName(e.target.value)}
                        /> */}
                        {/* <button onClick={handleDownload}>Download File</button> */}
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