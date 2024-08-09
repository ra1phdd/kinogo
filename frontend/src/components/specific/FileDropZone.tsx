import React, { useState, useRef } from 'react';

interface FileDropZoneProps {
    onFileUploaded: () => void;
}

const FileDropZone: React.FC<FileDropZoneProps> = ({ onFileUploaded }) => {
    const [uploadProgress, setUploadProgress] = useState(0);
    const [isDragActive, setIsDragActive] = useState(false);
    const [isFileLoading, setIsFileLoading] = useState(false);
    const [isFileUploaded, setIsFileUploaded] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleDragEnter = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setIsDragActive(true);
    };

    const handleDragLeave = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setIsDragActive(false);
    };

    const handleDragOver = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setIsDragActive(true);
    };

    const handleDrop = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setIsDragActive(false);

        if (event.dataTransfer.files && event.dataTransfer.files.length > 0) {
            setIsFileLoading(true);
            const files = event.dataTransfer.files;
            processingFile(files);
        }
    };

    const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files && event.target.files.length > 0) {
            setIsFileLoading(true);
            const files = event.target.files;
            processingFile(files);
        }
    };

    const handleClick = () => {
        fileInputRef.current?.click();
    };

    const processingFile = async (files: FileList) => {
        const formData = new FormData();
        for (let i = 0; i < files.length; i++) {
            formData.append('file', files[i]);
        }

        const xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://localhost:4000/upload', true);

        xhr.upload.onprogress = (event) => {
            if (event.lengthComputable) {
                const percentComplete = Math.round((event.loaded / event.total) * 100);
                setUploadProgress(percentComplete);
            }
        };

        xhr.onload = () => {
            if (xhr.status === 200) {
                setIsFileUploaded(true);
                onFileUploaded(); // Уведомление о завершении загрузки файла
            }
            setIsFileLoading(false);
        };

        xhr.onerror = () => {
            console.error('Upload failed');
            setIsFileLoading(false);
        };

        xhr.send(formData);
    };

    return (
        <div
            className="filedrop"
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            onDragOver={handleDragOver}
            onDrop={handleDrop}
            onClick={handleClick}
            style={{ cursor: 'pointer' }}
        >
            <input
                ref={fileInputRef}
                type="file"
                onChange={handleFileSelect}
                style={{ display: 'none' }}
            />
            {isDragActive ? (
                <p>Осталось только отпустить ПКМ</p>
            ) : isFileLoading ? (
                <p>Файл загружается... {uploadProgress}%</p>
            ) : isFileUploaded ? (
                <p>Файл загружен</p>
            ) : (
                <p>Перетащите файл сюда или нажмите для выбора файла</p>
            )}
        </div>
    );
};

export default FileDropZone;
