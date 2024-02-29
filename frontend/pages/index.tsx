import React, { useState, useEffect } from 'react';
import axios from 'axios';
import UploadForm from '../components/UploadForm';

interface UploadedFile {
    _id: string;
    courseName: string;
    batch: string;
    instructor: string;
    type: string;
    remark: string;
    link: string; // Add the link property
}

const IndexPage: React.FC = () => {
    const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[] | null>(null); // Initialize with null
    const [showUploadForm, setShowUploadForm] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [searchQuery, setSearchQuery] = useState('');
    const filesPerPage = 12; // Adjusted to show 9 tiles per page

    const fetchUploadedFiles = async () => {
        try {
            const response = await axios.get<UploadedFile[]>(`http://localhost:8080/api/fetch`);
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

    const paginate = (pageNumber: number) => setCurrentPage(pageNumber);

    const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchQuery(e.target.value);
        setCurrentPage(1); // Reset pagination to the first page when performing a new search
    };

    const filteredFiles = uploadedFiles?.filter(file =>
        file.courseName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.instructor.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.batch.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.remark.toLowerCase().includes(searchQuery.toLowerCase())
    ) || [];

    const indexOfLastFile = currentPage * filesPerPage;
    const indexOfFirstFile = indexOfLastFile - filesPerPage;
    const currentFiles = filteredFiles.slice(indexOfFirstFile, indexOfLastFile);

    const maxPages = Math.ceil(filteredFiles.length / filesPerPage);

    return (
        <div className="container mx-auto mt-8 relative">
            <h1 className="text-6xl py-4 font-bold text-center mb-4">STUDHELP-IITK</h1>
            <div className="mb-8 text-right">
                <input
                    type="text"
                    placeholder="Search..."
                    value={searchQuery}
                    onChange={handleSearch}
                    className="border my-4 text-black border-gray-300 rounded-full px-4 py-2 mx-auto block max-w-3xl w-full"
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
                {currentFiles.map((file, index) => (
                    <a href={file.link} download key={index}>
                        <div className="text-center mt-2 bg-gray-100 py-10 rounded-md shadow-md max-w-xs sm:max-w-full">
                            <h2 className="text-base text-gray-700 font-semibold">{file.courseName}</h2>
                            <p className="text-gray-500 text-xs">{file.type}</p>
                            <p className="text-gray-500 text-xs">Instructor: {file.instructor}, Batch: {file.batch}</p>
                            <p className="text-gray-500 text-xs">Remark: {file.remark}</p>
                        </div>
                    </a>
                ))}
            </div>

            {/* Pagination */}
            {filteredFiles.length > filesPerPage && (
                <div className="flex justify-center mt-8">
                    <button
                        onClick={() => paginate(currentPage - 1)}
                        disabled={currentPage === 1}
                        className="bg-gray-200 text-gray-800 font-semibold py-2 px-4 mx-1 rounded"
                    >
                        {"<<"}
                    </button>
                    {Array.from({ length: Math.min(maxPages, 10) }, (_, i) => {
                        const pageNumber = currentPage + i - 5;
                        return (
                            pageNumber > 0 && pageNumber <= maxPages && (
                                <button
                                    key={i}
                                    onClick={() => paginate(pageNumber)}
                                    className={`bg-gray-200 text-gray-800 font-semibold py-2 px-4 mx-1 rounded ${
                                        currentPage === pageNumber ? 'bg-gray-400' : ''
                                    }`}
                                >
                                    {pageNumber}
                                </button>
                            )
                        );
                    })}
                    <button
                        onClick={() => paginate(currentPage + 1)}
                        disabled={currentPage === maxPages}
                        className="bg-gray-200 text-gray-800 font-semibold py-2 px-4 mx-1 rounded"
                    >
                        {">>"}
                    </button>
                </div>
            )}


        </div>
    );
};

export default IndexPage;
