import { useEffect, useState } from "react";

import { ScrollHandler } from "utils/effects";
import { useFromTo } from "utils/hooks";
import { RandomKey } from "utils/content";
import Parsel from "components/parsel/parsel";
import Traveler from "components/traveler/traveler";

import styled from "styled-components";

const SHistory = styled.section`
    padding: 1rem;

    & .history-tabs {
        display: flex;
        align-items: center;

        & span {
            margin: .5rem;
            padding: 1rem;
            border-radius: 10px;
            transition: var(--transitionApp);
            cursor: pointer;

            &.active,
            &:hover {
                color: white;
                background: var(--blueColor);
            }
        }
    }

    & .history {
        padding: 1rem;
        max-height: 60vh;
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
    const [prevTab, setPrevTab] = useState();
    const { datalist, isStopLoad, isLoaded, getPart, zeroState, setDataList } = useFromTo([], 1)

    const [getType, getTypeOnRus, Item] = configHistoryParams(tab);

    const changeItem = (id, newData) => {
        const index = datalist.findIndex(d => d.id === id)
        datalist[index] = newData
        setDataList([...datalist]);
    }

    const removeItem = id => setDataList([...datalist.filter(d => d.id !== id)])

    useEffect(() => {
        if (prevTab !== tab) {
            zeroState();
        }
        if (datalist.length === 0 && !isLoaded) {
            loadHistory(getType, getTypeOnRus, getPart)
        }
        setPrevTab(tab)
    }, [datalist, isLoaded, getType, getTypeOnRus, prevTab, tab, getPart, zeroState]);

    return (
        <SHistory>
            <div className="history-tabs">
                <span className={tab === 0 ? 'active' : ''} onClick={() => setTab(0)}>Ваши посылки</span>
                <span className={tab === 1 ? 'active' : ''} onClick={() => setTab(1)}>Ваши путешествия</span>
            </div>

            {
                datalist.length > 0 &&
                <div className="history" onScroll={e => ScrollHandler(e, isStopLoad, false, () => loadHistory(getType, getTypeOnRus, getPart))}>
                    {datalist.map(d => <Item key={RandomKey()} data={d} isMy={true} changeItem={changeItem} removeItem={removeItem} />)}
                </div>
            }
        </SHistory>

    )
}