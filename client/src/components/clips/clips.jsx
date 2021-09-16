import { PreloadFile } from "utils/file";
import PreloadedFilesPlash from '../preloaded-files-plash/plash';

import styled from "styled-components";

const SClipItem = styled.div`
    width: 3rem;
    height: 3rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: ${props => props.color ? props.color : 'rgba(255, 255, 255, 0.38)'}; 
    border-radius: 50px;
    cursor: pointer;

    & img {
        width: 80%;
        height: 80%;
    }

    &:hover {
        background: var(--purpleColor);
    }
`;

const ClipOneBtn = ({ color, alt, srcIcon, onClick }) => {
    return (
        <SClipItem color={color} onClick={onClick}>
            <img src={srcIcon} alt={alt} />
        </SClipItem>
    )
}

export default function ClipPlash({ preloadedFiles = [], setFiles = () => { } }) {
    const addToPlash = (...files) => setFiles([...preloadedFiles, ...files]);
    const removeFile = filename => setFiles(preloadedFiles.filter(file => file.filename !== filename))

    const preloadedCB = (file, src, type) => {
        addToPlash({
            'type': type,
            'file': file,
            'src': src,
            'filename': file.name,
        });
    }

    return (
        <>
            <ClipOneBtn alt="clip" srcIcon="/assets/app/send-gallery.png" onClick={() => PreloadFile('image/*', preloadedCB)} />
            <PreloadedFilesPlash preloadedFiles={preloadedFiles} removeFile={removeFile} />
        </>
    )
}