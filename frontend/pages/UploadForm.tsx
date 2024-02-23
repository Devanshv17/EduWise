import React, { useState } from 'react';
import axios from 'axios';

interface UploadFormProps {
    fetchUploadedFiles: () => void;
    onClose: () => void; // Function to close the form
}

const UploadForm: React.FC<UploadFormProps> = ({ fetchUploadedFiles, onClose }) => {
    const [courseName, setCourseName] = useState('');
    const [batch, setBatch] = useState('');
    const [instructor, setInstructor] = useState('');
    const [type, setType] = useState('');
    const [detail, setDetail] = useState('');
    const [remark, setRemark] = useState('');
    const [photo, setPhoto] = useState<File | null>(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            setPhoto(e.target.files[0]);
        }
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append('courseName', courseName);
        formData.append('batch', batch);
        formData.append('instructor', instructor);
        formData.append('type', type);
        formData.append('detail', detail);
        formData.append('remark', remark);
        if (photo) {
            formData.append('photo', photo);
        }

        try {
            await axios.post('http://localhost:8080/api/upload', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
            alert('Upload successful!');
            fetchUploadedFiles();
        } catch (error) {
            alert('Error uploading data');
            console.error(error);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="bg-white shadow-md rounded-2xl px-8 pt-6 pb-8 mb-4 w-96 relative">
            <button onClick={onClose} className="absolute top-2 right-4 mt-2 mr-2 text-gray-500 hover:text-gray-700">
                Close
            </button>
            <h1 className="text-3xl text-blue-700 font-semibold mb-6">Upload Course Information</h1>
            <div className="space-y-4">
                <label className="block text-gray-500 mb-2">Course Name:</label>
                <input
                    type="text"
                    value={courseName}
                    onChange={(e) => setCourseName(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500 text-black"
                    placeholder="Enter course name"
                />
                <label className="block text-gray-500 mb-2">Batch:</label>
                <input
                    type="text"
                    value={batch}
                    onChange={(e) => setBatch(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500 text-black"
                    placeholder="Enter batch"
                />
                <label className="block text-gray-500 mb-2">Instructor:</label>
                <input
                    type="text"
                    value={instructor}
                    onChange={(e) => setInstructor(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500 text-black"
                    placeholder="Enter instructor"
                />
                <label className="block text-gray-500 mb-2">Type:</label>
                <select
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500 appearance-none text-black"
                >
                    <option value="">Select Type</option>
                    <option value="Question Paper">Question Paper</option>
                    <option value="Notes">Notes</option>
                </select>
                <label className="block text-gray-500 mb-2">Detail:</label>
                <select
                    value={detail}
                    onChange={(e) => setDetail(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500 appearance-none text-black"
                >
                    <option value="">Select Detail</option>
                    <option value="Lecture Number">Lecture Number</option>
                    <option value="Midsem">Midsem</option>
                    <option value="Endsem">Endsem</option>
                    <option value="Quiz">Quiz</option>
                    {/* Add more options as needed */}
                </select>
                <label className="block text-gray-500 mb-2">Remark:</label>
                <input
                    type="text"
                    value={remark}
                    onChange={(e) => setRemark(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500 text-black"
                    placeholder="Enter remark"
                />
                <label className="block text-gray-500 mb-2">Photo:</label>
                <input type="file" accept="image/*" onChange={handleFileChange} className="w-full focus:outline-none" />
            </div>
            <button type="submit" className="w-full bg-blue-500 text-white py-2 my-4 px-4 rounded-md hover:bg-blue-600 transition duration-300">Upload</button>
        </form>
    );
};

export default UploadForm;
