import React, { useState, useEffect } from 'react';
import axios from 'axios';
import UploadForm from './UploadForm';

interface UploadedFile {
    courseName: string;
    batch: string;
    photo: string;
    instructor: string;
    type: string;
    remark: string;
}

const IndexPage: React.FC = () => {
    const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
    const [showUploadForm, setShowUploadForm] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [searchQuery, setSearchQuery] = useState('');
    const filesPerPage = 6; // Adjusted to show 9 tiles per page

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
    }, []);

    const handleToggleForm = () => {
        setShowUploadForm(!showUploadForm);
    };

    const indexOfLastFile = currentPage * filesPerPage;
    const indexOfFirstFile = indexOfLastFile - filesPerPage;
    const currentFiles = uploadedFiles.slice(indexOfFirstFile, indexOfLastFile);

    const paginate = (pageNumber: number) => setCurrentPage(pageNumber);

    const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchQuery(e.target.value);
    };

    const filteredFiles = currentFiles.filter(file =>
        file.courseName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.instructor.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.batch.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.remark.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="container mx-auto mt-8 relative">
            <h1 className="text-6xl py-4 font-bold text-center mb-4">STUDHELP-IITK</h1> {/* Heading added here */}
            <div className="mb-8 text-right">
            <input
                type="text"
                placeholder="Search..."
                value={searchQuery}
                onChange={handleSearch}
                className="border text-black border-gray-300 rounded-full px-4 py-2 mx-auto block max-w-3xl w-full"
            />

                <button onClick={handleToggleForm} className="bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 transition duration-300">
                    {showUploadForm ? 'Close Form' : 'Upload Files'}
                </button>
            </div>

            {showUploadForm && (
                <div className="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50 bg-white p-8 rounded-2xl shadow-lg">
                    <UploadForm fetchUploadedFiles={fetchUploadedFiles} onClose={() => setShowUploadForm(false)} />
                </div>
            )}

            <div className="grid gap-16 mt-2 grid-cols-1 sm:grid-cols-2 md:grid-cols-3">
                {filteredFiles.map((file, index) => (
                    <div key={index} className="bg-gray-100 p-2 rounded-md shadow-md max-w-xs sm:max-w-full">
                        <img src={`/uploads/${file.photo}`} alt={file.courseName} className="w-full h-auto" />
                        <div className="text-center mt-2">
                            <h2 className="text-base text-gray-700 font-semibold">{file.courseName}</h2>
                            <p className="text-gray-500 text-xs">{file.type}</p>
                            <p className="text-gray-500 text-xs">Instructor: {file.instructor}, Batch: {file.batch}</p>
                            <p className="text-gray-500 text-xs">Remark: {file.remark}</p>
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
