// IndexPage.tsx

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import UploadForm from './UploadForm';

interface UploadedFile {
  courseName: string;
  year: string;
  photo: string;
}

const IndexPage: React.FC = () => {
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [filesPerPage] = useState(12); // 3 tiles horizontally and 4 tiles vertically

  const fetchUploadedFiles = async () => {
    try {
      const response = await axios.get<UploadedFile[]>('http://localhost:8080/api/fetch');
      setUploadedFiles(response.data);
    } catch (error) {
      console.error('Error fetching uploaded files:', error);
    }
  };

  useEffect(() => {
    fetchUploadedFiles();
  }, []); // Fetch uploaded files on component mount

  const handleToggleForm = () => {
    setShowUploadForm(!showUploadForm);
  };

  // Pagination
  const indexOfLastFile = currentPage * filesPerPage;
  const indexOfFirstFile = indexOfLastFile - filesPerPage;
  const currentFiles = uploadedFiles.slice(indexOfFirstFile, indexOfLastFile);

  const paginate = (pageNumber: number) => setCurrentPage(pageNumber);

  return (
    <div className="container mx-auto mt-8 relative">
      <button onClick={handleToggleForm} className="bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 transition duration-300 absolute top-0 right-0 mt-4 mr-4">
        {showUploadForm ? 'Close Form' : 'Upload Files'}
      </button>

      {showUploadForm && <UploadForm fetchUploadedFiles={fetchUploadedFiles} />}

      <div className="grid grid-cols-3 gap-4 mt-10">
        {currentFiles.map((file, index) => (
          <div key={index} className="bg-gray-100 p-4 rounded-md shadow-md">
            <img src={`/uploads/${file.photo}`} alt={file.courseName} className="w-full h-auto" />
            <div className="text-center mt-4">
              <h2 className="text-lg font-semibold">{file.courseName}</h2>
              <p className="text-gray-500">{file.year}</p>
            </div>
          </div>
        ))}
      </div>
      {/* Pagination */}
      <ul className="flex justify-center mt-8">
        {Array.from({ length: Math.ceil(uploadedFiles.length / filesPerPage) }, (_, i) => (
          <li key={i}>
            <button
              onClick={() => paginate(i + 1)}
              className={`bg-gray-200 text-gray-800 font-semibold py-2 px-4 mx-1 rounded ${
                currentPage === i + 1 ? 'bg-gray-400' : ''
              }`}
            >
              {i + 1}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default IndexPage;
