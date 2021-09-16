import { useCallback } from "react";

import { GetValueFromListByIDAndInputValue, OnChangeTransitPoint, ScrollHandler } from "utils/effects";
import { useInput } from "utils/form";
import { RandomKey } from "utils/content";
import { useFromTo } from "utils/hooks";
import Input from "components/form-input/input";
import Parsel from "components/parsel/parsel";

import styled from "styled-components";

const SParsels = styled.section`
    padding: 1rem;
    background: var(--blueColor);

    & .filters {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        justify-content: space-evenly;

        & * {
            color: var(--onHoverColor);
        }

        & .search_btn {
            padding: .5rem 1rem;
            background: var(--darkGreyColor);
            border-radius: 10px;
            cursor: pointer;
            box-shadow: var(--boxShadow);
            transition: var(--transitionApp);

            &:hover {
                background: var(--onHoverBG);
            }
        }
    }

    @media screen and (max-width: 600px) {
        & .filters {
            justify-content: start;
        }
    }
`;


export default function ParselsPage() {
    const from = useInput('');
    const to = useInput('');
    const fromID = useInput('');
    const toID = useInput('');
    const startDT = useInput('');
    const endDT = useInput('');

    from.base.onChange = e => OnChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => OnChangeTransitPoint(to, e, toID.setCertainValue);

    const { datalist, isStopLoad, getPart } = useFromTo()

    const loadParsels = useCallback(() => getPart("parsels", {
        'from': GetValueFromListByIDAndInputValue("from-list", from.base.value),
        'to':  GetValueFromListByIDAndInputValue("to-list", to.base.value),
        'startDT': Date.parse(startDT.base.value),
        'endDT': Date.parse(endDT.base.value)
    }, 'Не удалось загрузить посылки'), [from, to, startDT, endDT, getPart])

    return (
        <SParsels>
            <div className="filters">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>

                <Input type="date" name="startDT" base={startDT.base} required={false} labelText="С:" />
                <Input type="date" name="endDT" base={endDT.base} required={false} labelText="До:" />

                <span className="search_btn" onClick={loadParsels}>
                    <i className="fa fa-search" aria-hidden="true"></i>
                </span>
            </div>

            {
                datalist &&
                <div className="parsels" onScroll={e => ScrollHandler(e, isStopLoad, false, loadParsels)}>
                    {datalist?.map(p => <Parsel key={RandomKey()} data={p} />)}
                </div>
            }
        </SParsels>
    )
}