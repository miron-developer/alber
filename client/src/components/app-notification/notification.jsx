import { useState } from 'react';

import { RandomKey } from 'utils/content';

import styled from 'styled-components';

const SAppNotifications = styled.div`
    position: fixed;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    z-index: 20;
`

const SAppNotification = styled.div`
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: max-content;
    max-width: 50vw;
    margin: .5rem;
    padding: 1rem;
    color: white;
    border-radius: 5px;
    background: #596d8c;
    box-shadow: var(--boxShadow);
`;

const SCloseNotification = styled.div`
    padding: 5px;
    margin-left: 10px;
    height: 1.5rem;
    width: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--redColor);
    background: var(--darkRedColor);
    border-radius: 50%;
    cursor: pointer;
`;

let add = () => {};
let removeNotification = id => {};
let interID;

export const Notify = (type, content, isAutoClose) => {
    if (add.length !== 0) return add(type, content, isAutoClose);
    interID = setInterval(() => {
        if (add.length !== 0) {
            clearInterval(interID);
            add(type, content, isAutoClose);
        }
    }, 0);
}

const Notification = ({ id = 0, type = "fail", content, isAutoClose = true }) => {
    if (!content) return <div></div>;
    if (isAutoClose) setTimeout(() => removeNotification(id), 5000);

    let noteContent = content
    if (content instanceof Function) noteContent = content(); 

    return (
        <SAppNotification>
            <div className="notification-content">{noteContent}</div>
            <SCloseNotification onClick={()=>removeNotification(id)}>
                <i className="fa fa-times" aria-hidden="true"></i>
            </SCloseNotification>
        </SAppNotification>
    )
}

export default function Notifications() {
    const [ntfs, setNTFS] = useState([]);

    removeNotification = id => setNTFS(ntfs.filter(item => item.props.id !== id));
    add = (type, content, isAutoClose)  => {
        const key = RandomKey();
        setNTFS([...ntfs, <Notification key={key} id={key} type={type} content={content} isAutoClose={isAutoClose} />]);
    }

    return <SAppNotifications>{ntfs}</SAppNotifications>;
}