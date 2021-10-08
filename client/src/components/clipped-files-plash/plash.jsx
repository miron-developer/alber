import { RandomKey } from "utils/content";
import { PopupOpen } from "../popup/popup";

import styled from "styled-components";

const SFilesPlash = styled.div`
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    background: #000000ad;
`;

const SClippedFileWrapper = styled.div`
    position: relative;
    max-width:  ${props => props.size ? props.size : '15rem'};
    width: max-content;
    max-height: ${props => props.size ? props.size : '15rem'};
    margin: .5rem;
    padding: .5rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: 'space-between';
    color: var(--onHoverColor);
    /* background: #000000ad; */
    border-radius: 10px;
    cursor: pointer;
    overflow: hidden;

    & span {
        position: absolute;
        bottom: 0;
        left: 0;
        padding: 5px;
        background: #0f6dfb9e;
    }
`;

const SClippedFileSrc = styled.div`
    width: 100%;
    
    & > * {
        max-width: 100%;
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


const RenderClippedFile = ({ id, filename, src, size, onClick = () => { }, removeFile }) => {
    return (
        <SClippedFileWrapper size={size} onClick={onClick} >
            <SRemoveFile onClick={e => e.stopPropagation() || removeFile(id, src)}>
                <i className="fa fa-times" aria-hidden="true"></i>
            </SRemoveFile>

            <SClippedFileSrc>
                <img src={src} alt="uploaded img" />
            </SClippedFileSrc>
            <span>{filename}</span>
        </SClippedFileWrapper>
    )
}

export default function ClippedFiles({ files = [], removeFile}) {
    return (
        <SFilesPlash>
            {
                files?.map(
                    file => <RenderClippedFile key={RandomKey()} {...file} 
                        onClick={e => e.preventDefault() || PopupOpen(RenderClippedFile, { ...file, 'size': '100%' })} 
                        removeFile={removeFile}
                    />
                )
            }
        </SFilesPlash>
    )
}