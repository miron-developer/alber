import { useCallback, useState } from "react";

import { OnChangeTransitPoint } from "utils/effects";
import { GetDataByCrieteries } from "utils/api";
import { useInput } from "utils/form";
import { RandomKey } from "utils/content";
import { Notify } from "components/app-notification/notification";
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

    const [parsels, setParsels] = useState();

    const getParsels = useCallback(async () => {
        const res = await GetDataByCrieteries('parsels', {
            'from': fromID.base.value,
            'to': toID.base.value,
            'startDT': startDT.base.value,
            'endDT': endDT.base.value
        });
        if (res.err && res.err !== "ok") return Notify('fail', "Посылок не найдено");
        setParsels(res)
    }, [fromID, toID, startDT, endDT])

    return (
        <SParsels>
            <div className="filters">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>

                <Input type="date" name="startDT" base={startDT.base} labelText="С:" />
                <Input type="date" name="endDT" base={endDT.base} labelText="До:" />

                <span className="search_btn" onClick={getParsels}>
                    <i className="fa fa-search" aria-hidden="true"></i>
                </span>
            </div>

            {
                parsels &&
                <div className="parsels">
                    {parsels?.map(p => <Parsel key={RandomKey()} data={p} />)}
                </div>
            }
        </SParsels>
    )
}