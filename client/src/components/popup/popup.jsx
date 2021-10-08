import { useState } from 'react';

import styled, { css } from 'styled-components';

const SDisplayFlexCenter = css`
    display: flex;
    align-items: center;
    justify-content: center;
`;

const SPopup = styled.div`
    position: fixed;
    left: 50%;
    top: 50%;
    width: 0;
    height: 0;
    ${SDisplayFlexCenter};
    border-radius: 5px;
    transition: var(--transitionApp);
    opacity: 0;
    z-index: -10;

    &.opened {
        left: 0;
        top: 0;
        width: 100vw;
        height: 100vh;
        opacity: 1;
        z-index: 20;
    }
`;


const SPopupBG = styled.div`
    position: absolute;
    width: 100vw;
    height: 100vh;
    background: #000000ab;
    backdrop-filter: blur(1px);
`;

const SPopupClose = styled.div`
    position: absolute;
    right: 20px;
    top: 10px;
    width: 1rem;
    height: 1rem;
    padding: 5px;
    ${SDisplayFlexCenter};
    background: var(--redColor);
    border-radius: 100%;
    transition: var(--transitionApp);
    cursor: pointer;

    &:hover{
        background: var(--darkRedColor);
    }
`;

const SContentWrapper = styled.div`
    max-height: 90vh;
    max-width: 80vw;
    height: max-content;
    width: max-content;
    ${SDisplayFlexCenter};
    flex-direction: column;
    border-radius: 10px;
    background: #ffffffc7;
    z-index: 20;

    & > * {
        overflow: auto;
    }
`;

let add = () => { };

let close = () => { };
let closeCB = () => { };

let onCloseCB = () => { };

export const OnClosePopup = (f = () => { }) =>  onCloseCB = f;

export const PopupOpen = (Component, componentArgs = {}) => {
    add(Component, componentArgs);
    if (componentArgs.closeCB) closeCB = componentArgs.closeCB;
}
export const ClosePopup = () => close();

export default function Popup() {
    const [popup, setPopup] = useState(null);
    const [isOpened, setOpenedState] = useState(false);
    add = (Component, componentArgs) => setPopup(<Component {...componentArgs} />) || setOpenedState(true);
    close = () => {
        setOpenedState(false) || setPopup(null);
        closeCB();
        onCloseCB()
    }

    return (
        <SPopup className={`${isOpened ? 'opened' : ''}`} >
            <SPopupBG onClick={close} />

            <SPopupClose onClick={close}>
                <i className="fa fa-times" aria-hidden="true"></i>
            </SPopupClose>

            <SContentWrapper>
                {popup}
            </SContentWrapper>
        </SPopup>
    )
}