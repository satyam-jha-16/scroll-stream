import axios from 'axios';
import React, { useState, useRef } from 'react';
import VideoPlayer from './VideoPlayer';

const VideoUploader = () => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploadStatus, setUploadStatus] = useState('');
  const [videoUrl, setVideoUrl] = useState('');
  

  const handleFileChange = (event) => {
    setSelectedFile(event.target.files[0]);
  };

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append('file', selectedFile);

    try {
      setUploadStatus('Uploading...');
      const response = await axios.post('http://localhost:8000/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });

      setUploadStatus('Upload complete!');
      setVideoUrl(response.data.videoUrl); 
    } catch (error) {
      setUploadStatus('Upload failed. Please try again.');
      console.error(error);
    }
  };

  const playerRef = useRef(null);
  // const videoLink =  ;
  const videoPlayerOptions = {
    controls: true,
    responsive: true,
    fluid: true,
    sources: [
      {
        src: videoUrl,
        type: "application/x-mpegURL",
      },
    ],
  };
  const handlePlayerReady = (player) => {
    playerRef.current = player;

    // You can handle player events here, for example:
    player.on("waiting", () => {
      videojs.log("player is waiting");
    });

    player.on("dispose", () => {
      videojs.log("player will dispose");
    });
  };

  return (
    <div>
      <input type="file" onChange={handleFileChange} />
      <button onClick={handleUpload} disabled={!selectedFile}>Upload Video</button>
      <p>{uploadStatus}</p>
      {videoUrl && (
        <VideoPlayer
        options={videoPlayerOptions}
        onReady={handlePlayerReady}
        /> 
      )}
    </div>
  );
};

export default VideoUploader;
