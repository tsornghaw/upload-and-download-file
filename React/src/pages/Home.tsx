import React, { ChangeEvent, ChangeEventHandler, useEffect, useState } from 'react';
import axios from 'axios';

// Define a type for your data
type DataItem = {
    Dataname: string;
    DownloadURL: string;
  };

const Home = (props: { name: string }) => {
    const [fileList, setFileList] = useState<FileList | null>(null);
    const [selectedImage, setSelectedImage] = useState('');

    const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
        setFileList(e.target.files);
        
        const selectedFiles = e.target.files

        if (selectedFiles && selectedFiles.length > 0) {
            const file = selectedFiles[0];
            setSelectedImage(URL.createObjectURL(file));
        }
    };

    const handleDownload = async () => {
        // Implement the API call to download here
        console.log('Downloading...');
        // You can use fetch or any other library to make the API call
        try {
            const response = await fetch('http://localhost:8000/api/auth/download/test.png');
            const blob = await response.blob();
            
            // Create a URL for the blob data and trigger a download
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'filename.extension'; // Set the desired filename
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Error downloading the file:', error);
        }
    };

    const handleApi = () => {
        if (fileList === null || fileList.length === 0) {
            return
        }
        const formData = new FormData()
        formData.append(`file`, fileList[0], fileList[0].name)
        console.log(fileList[0])
        fetch('http://localhost:8000/api/auth/upload', {
            method: 'POST',
            body: formData,
        })
    }




    // Testing
    // const [data, setData] = useState<DataItem[]>([]);
    // const [selectedOption, setSelectedOption] = useState<string>('');

    // useEffect(() => {
    //     // Make an HTTP request to your Go backend
    //     axios.get('http://localhost:8000/api/auth/UserSearchAllData')
    //         .then(response => {
    //             console.log(response.data)
    //             setData(response.data);
    //         })
    //         .catch(error => {
    //             // Handle the error
    //         });
    // }, []);
    // Define a state variable to store the options fetched from the API
    const [options, setOptions] = useState<DataItem[]>([]);

    // Define a state variable to store the selected option
    const [selectedOption, setSelectedOption] = useState('');
    // Use the useEffect hook to fetch data from the API when the component mounts
    useEffect(() => {
        // Replace 'your-api-endpoint' with the actual API endpoint URL
        fetch('your-api-endpoint')
        .then((response) => response.json())
        .then((data) => {
            // Assuming the API response is an array of objects with 'value' and 'label' properties
            setOptions(data);
        })
        .catch((error) => {
            console.error('Error fetching data from the API:', error);
        });
    }, []); // The empty dependency array ensures this effect runs only once on component mount

    // Handle the selection change
    const handleSelectChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
        setSelectedOption(event.target.value);
    };

    
    return (
        <div>
            {props.name ? 
                (
                    <div>
                        // If Authenticated
                        <input type="file" id="myFile" name="filename" accept="image/*" onChange={handleFileChange}/>
                        {selectedImage && (
                            <img
                            src={selectedImage}
                            alt="Selected"
                            style={{ maxWidth: '600px', maxHeight: '600px' }} // Set your desired max width and max height
                            />
                        )}
                        <br />
                        {/* <input type="submit" /> */}
                        <button onClick={handleApi}>Submit</button>
                        <br />
                        <button onClick={handleDownload}>Download File</button>

                        {/* Testing */}
                        {/* <select
                            value={selectedOption}
                            onChange={(e) => setSelectedOption(e.target.value)}
                        >
                            <option value="">Select an option</option>
                            {data.map((item) => (
                            <option key={item.Dataname} value={item.DownloadURL}>
                                {item.Dataname}
                            </option>
                            ))}
                        </select> */}
                        <br />
                        <label>Select an option: </label>
                        <select value={selectedOption} onChange={handleSelectChange}>
                            <option value="">Select an option</option>
                            {options.map((option) => (
                            <option key={option.Dataname} value={option.Dataname}>
                                {option.Dataname}
                            </option>
                            ))}
                        </select>
                        {selectedOption && <p>Selected option: {selectedOption}</p>}
                    </div>
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