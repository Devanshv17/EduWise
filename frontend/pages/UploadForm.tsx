// UploadForm.tsx

import React, { useState } from 'react';
import axios from 'axios';

interface UploadFormProps {
  fetchUploadedFiles: () => void;
}

const UploadForm: React.FC<UploadFormProps> = ({ fetchUploadedFiles }) => {
  const [courseName, setCourseName] = useState('');
  const [year, setYear] = useState('');
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
    formData.append('year', year);
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
      // After successful upload, fetch the updated list of files
      fetchUploadedFiles();
    } catch (error) {
      alert('Error uploading data');
      console.error(error);
    }
  };

  return (
    <div className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4 w-96">
      <h1 className="text-3xl font-semibold mb-6">Upload Course Information</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block mb-2">Course Name:</label>
          <input
            type="text"
            value={courseName}
            onChange={(e) => setCourseName(e.target.value)}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500"
            placeholder="Enter course name"
          />
        </div>
        <div>
          <label className="block mb-2">Year:</label>
          <input
            type="text"
            value={year}
            onChange={(e) => setYear(e.target.value)}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500"
            placeholder="Enter year"
          />
        </div>
        <div>
          <label className="block mb-2">Photo:</label>
          <input type="file" accept="image/*" onChange={handleFileChange} className="w-full focus:outline-none" />
        </div>
        <button type="submit" className="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 transition duration-300">Upload</button>
      </form>
    </div>
  );
};

export default UploadForm;
