import { RandomKey } from "utils/content";

import styled from "styled-components";

const SFilesPlash = styled.div`
    position: absolute;
    left: 0;
    right: 0;
    top: 100%;
    background: #000000ad;
    overflow: auto;
    white-space: nowrap;
`;

const SFileUploadWrapper = styled.div`
    position: relative;
    width:  4rem;
    display: inline-block;
    margin: .5rem;
    padding: .5rem;
    overflow: hidden;

    & > * {
        width: 100%;
        color: var(--onHoverColor);
    }
`;

const SRemoveFile = styled.div`
    position: absolute;
    right: 0;
    top: 0;
    width: 1rem;
    height: 1rem;
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--redColor);
    border-radius: 100%;
    transition: var(--transitionApp);
    cursor: pointer;

    &:hover{
        background: var(--darkRedColor);
    }
`;

const RenderUploadedFile = ({filename, src, removeFile}) => {
    return (
        <SFileUploadWrapper>
            <SRemoveFile onClick={()=>removeFile(filename)}>
                <i className="fa fa-times" aria-hidden="true"></i>
            </SRemoveFile>

            <img src={src} alt="uploaded img" />
            <div>{filename}</div>
        </SFileUploadWrapper>
    )
}

export default function PreloadedFilesPlash({ preloadedFiles, removeFile = ()=>{}}) {
    return (
        <SFilesPlash>
            {
                preloadedFiles.map(
                    file => <RenderUploadedFile key={RandomKey()} {...file} removeFile={filename => removeFile(filename)} />
                )
            }
        </SFilesPlash>
    )
}