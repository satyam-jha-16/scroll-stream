import Upload from "./components/Upload.jsx"
import New from "./components/New.jsx"
import './App.css'
import React, { useState } from 'react';
import ReactPlayer from 'react-player';
import More from "./components/More.jsx"

function App() {

  return (
    <>
    <More /> 
      {/* <Upload /> */}
      {/* <New /> */}\
      {/* <ReactPlayer
                url= 'https://firebasestorage.googleapis.com/v0/b/streamlit-f6eda.appspot.com/o/videos%2F009b2292-7d71-4f2f-bef2-12fd5e0493f7%2Fnew%2Fsegment000.ts?alt=media&token=e09d83f7-d202-46c1-a8a4-95e92c696d6d'
                controls
                width="100%"
                height="auto"
                config={{
                    file: {
                        forceHLS: true,
                    },
                }}
            /> */}
    </>
  )
}

export default App
