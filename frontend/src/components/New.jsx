import React, { useState } from 'react';
import ReactPlayer from 'react-player';

const VideoUploader = () => {
    const [videoUrl, setVideoUrl] = useState('');

    const uploadVideo = async (event) => {
        const file = event.target.files[0];
        const formData = new FormData();
        formData.append('file', file);

        const response = await fetch('http://localhost:8000/upload', {
            method: 'POST',
            body: formData,
        });

        const data = await response.json();
        setVideoUrl(data.videoUrl);
    };

    return (
        <div>
            <img src='https://firebasestorage.googleapis.com/v0/b/streamlit-f6eda.appspot.com/o/videos%2F009b2292-7d71-4f2f-bef2-12fd5e0493f7%2F3373670.jpg?alt=media&token=3417c10d-8341-4fe4-8587-663721b35206' />
            <input type="file" onChange={uploadVideo} />
            {videoUrl && (
                <ReactPlayer
                    url={videoUrl}
                    controls
                    config={{
                        file: {
                            attributes: {
                                crossOrigin: 'anonymous',
                            },
                            forceHLS: true,
                        },
                    }}
                />
            )}
        </div>
    );
};

export default VideoUploader;
