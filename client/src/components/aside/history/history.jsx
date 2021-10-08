import { useEffect, useState } from "react";

import { ScrollHandler } from "utils/effects";
import { useFromTo } from "utils/hooks";
import { RandomKey } from "utils/content";
import Parsel from "components/parsel/parsel";
import Traveler from "components/traveler/traveler";

import styled from "styled-components";

const SHistory = styled.section`
    padding: 1rem 0;
    height: 100%;

    & .history-tabs {
        display: flex;
        align-items: center;

        & span {
            margin: .5rem;
            padding: 1rem;
            color: var(--onHoverColor);
            border-radius: 10px;
            transition: var(--transitionApp);
            cursor: pointer;

            &.active,
            &:hover {
                background: var(--blueColor);
            }
        }
    }

    & .history {
        height: 70%;
        overflow: auto;
        border-radius: 10px;
        background: var(--offHoverBG);
    }
`

const loadHistory = (getType, getTypeOnRus, getPart) => getPart(getType, { 'type': 'user' }, 'Не удалось загрузить ' + getTypeOnRus)

const configHistoryParams = tab => {
    if (tab === 0) return ['parsels', 'посылки', Parsel];
    return ['travelers', 'путешествия', Traveler]
}

export default function History() {
    const [tab, setTab] = useState(0);
    const [isReload, setReload] = useState(0);
    const { datalist, isStopLoad, isLoaded, getPart, zeroState, setDataList } = useFromTo()

    const [getType, getTypeOnRus, Item] = configHistoryParams(tab);

    const changeItem = (id, newData) => {
        const index = datalist.findIndex(d => d.id === id)
        datalist[index] = newData
        setDataList([...datalist]);
    }

    const reloadItem = () => {
        setReload(true);
        setTimeout(() => {
            setReload(false);    
        }, 50);
    }

    const removeItem = id => setDataList([...datalist.filter(d => d.id !== id)]);


    useEffect(() => {
        if (datalist.length === 0 && !isLoaded) {
            loadHistory(getType, getTypeOnRus, getPart)
        }
    }, [datalist, isLoaded, getType, getTypeOnRus, tab, getPart, zeroState]);

    return (
        <SHistory>
            <div className="history-tabs">
                <span className={tab === 0 ? 'active' : ''} onClick={() => setTab(0) || zeroState()}>Ваши посылки</span>
                <span className={tab === 1 ? 'active' : ''} onClick={() => setTab(1) || zeroState()}>Ваши путешествия</span>
            </div>

            {
                datalist.length > 0
                    ? <div className="history" onScroll={e => ScrollHandler(e, isStopLoad, false, () => loadHistory(getType, getTypeOnRus, getPart))}>
                        {datalist.map(d => <Item key={RandomKey()} data={d} isMy={true} changeItem={changeItem} removeItem={removeItem} reloadItem={reloadItem} />)}
                    </div>
                    : <div className="history">Отсутствует</div>
            }
        </SHistory>

    )
}