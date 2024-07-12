import axios from 'axios';
import React, { useState, useRef } from 'react';
import VideoPlayer from './VideoPlayer';

const More = () => {
  const videoUrl = 'https://firebasestorage.googleapis.com/v0/b/streamlit-f6eda.appspot.com/o/videos%2Fnew%2Findex.m3u8?alt=media&token=f4ddb509-9fab-4e2c-b673-88734fa6a4c6'

  const playerRef = useRef(null);
  // const videoLink =  ;
  const videoPlayerOptions = {
    controls: true,
    responsive: true,
    fluid: true,
    sources: [
      {
        src: "https://storage.googleapis.com/streamlit--bucket/1af26ca0-d556-4782-b675-d4db532398d8/index.m3u8",
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


      <VideoPlayer
        options={videoPlayerOptions}
        onReady={handlePlayerReady}
      />
    </div>
  );
};

export default More;
