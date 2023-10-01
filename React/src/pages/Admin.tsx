import React, { ChangeEvent, ChangeEventHandler, useEffect, useState } from 'react';

type DataItems = {
    id: Uint16Array;
    upload_time: string;
    share_time: string;
    share_limit: number;
    file_size: number;
    file_name: string;
    file_type: string;
    file_content: string;
    download_url: string;
}

const Admin = (props: { name: string }) => {
    // Define a state variable to store the options fetched from the API
    const [fileList, setFileList] = useState<FileList | null>(null);
    const [selectedImage, setSelectedImage] = useState('');
    const [options, setOptions] = useState<DataItems[]>([]);
    const [selectedOption, setSelectedOption] = useState<string>("");
    const [downloadmessage, setdownloadMessage] = useState('');
    const [uploadmessage, setuploadMessage] = useState('');
    const [deleteResponse, setDeleteResponse] = useState('');
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
        fetch('http://localhost:8000/api/auth/admin/SearchAllData', {
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

    const handleDelete = (downloadUrl: string) => {
        // Make an API call to delete data using downloadUrl
        fetch(`http://localhost:8000/api/auth/admin/deletedata`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ download_url: downloadUrl }),
            credentials: 'include',
        })
        .then((response) => {
            if (response.status === 200) {
                setDeleteResponse('Data deleted successfully');
            } else {
                setDeleteResponse('Failed to delete data');
            }
        })
        .catch((error) => {
            console.error('Error deleting data:', error);
        });

        // Re-fetch data and update options when the button is clicked
        fetch('http://localhost:8000/api/auth/admin/SearchAllData', {
            credentials: 'include',
        })
        .then((response) => response.json())
        .then((data) => {
            setOptions(data);
        })
        .catch((error) => {
            console.error('Error fetching data from the API:', error);
        });
    };

    // Use the useEffect hook to fetch data from the API when the component mounts
    useEffect(() => {
        fetch('http://localhost:8000/api/auth/admin/SearchAllData', {
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
        setSelectedOption(event.target.value);
    };

    const handleRangeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        // Update the state with the current range value
        setDownloadTimes(parseInt(event.target.value, 10));
    };
    
    const tableStyle: React.CSSProperties = {
        width: '100%', // Set the table width to 100%
        borderCollapse: 'collapse', // Merge adjacent cell borders
        border: '1px solid #ccc', // Add a border to the table
    };

    const thStyle: React.CSSProperties = {
        backgroundColor: '#f2f2f2',
        padding: '10px',
        border: '1px solid #000', // Add a border to table header cells
    };
    
    const tdStyle: React.CSSProperties = {
        padding: '8px',
        border: '1px solid #ccc',
    };

    const buttonStyleadmin: React.CSSProperties = {
        backgroundColor: '#007bff', // Blue background color
        color: '#fff', // White text color
        padding: '10px 20px', // Padding around the text
        border: 'none', // Remove the default button border
        borderRadius: '4px', // Add rounded corners
        cursor: 'pointer', // Change cursor to a pointer on hover
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
                            style={{ maxWidth: '600px', maxHeight: '600px' }}
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
                            <option key={option.file_name + index} value={option.download_url}>
                            {option.file_name}
                            </option>
                        ))}
                        </select>
                        {selectedOption && <p>Selected option: {selectedOption}</p>}
                        <button onClick={handleDownload}>Download File</button>
                        <p>{downloadmessage}</p>
                        <br  />

                        {/* If Authenticated - Data List and Delete */}
                        <table style={tableStyle}>
                            <thead>
                            <tr>
                                <th style={thStyle}>File_name</th>
                                <th style={thStyle}>Upload_time</th>
                                <th style={thStyle}>Share_limit</th>
                                <th style={thStyle}>File_size</th>
                                <th style={thStyle}>File_type</th>
                                <th style={thStyle}>download_url</th>
                                <th style={thStyle}>edit_button</th>
                            </tr>
                            </thead>
                            <tbody>
                                {options && options.map((option) => (
                                    <tr key={option.file_name}>
                                        <td style={tdStyle}>{option.file_name}</td>
                                        <td style={tdStyle}>{option.upload_time}</td>
                                        <td style={tdStyle}>{option.share_limit}</td>
                                        <td style={tdStyle}>{(option.file_size / 1000).toFixed(2)} KB</td>
                                        <td style={tdStyle}>{option.file_type}</td>
                                        <td style={tdStyle}>"http://localhost:8000/" + {option.download_url}</td>
                                        <td style={tdStyle}>
                                            <button style={buttonStyleadmin} onClick={() => handleDelete(option.download_url)}>
                                                Delete
                                            </button>
                                            {deleteResponse}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
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

export default Admin;


function setSelectedImage(arg0: string) {
    throw new Error('Function not implemented.');
}
// Optimization: pjchender - Day23