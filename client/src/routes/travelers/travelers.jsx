import { useCallback, useState } from "react";

import { OnChangeTransitPoint } from "utils/effects";
import { GetDataByCrieteries } from "utils/api";
import { useInput } from "utils/form";
import { RandomKey } from "utils/content";
import { Notify } from "components/app-notification/notification";
import Input from "components/form-input/input";
import Traveler from "components/traveler/traveler";

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

export default function TravelersPage() {
    const from = useInput('');
    const to = useInput('');
    const fromID = useInput('');
    const toID = useInput('');
    const departure = useInput('');
    const arrival = useInput('');

    from.base.onChange = e => OnChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => OnChangeTransitPoint(to, e, toID.setCertainValue);

    const [travelers, setTravelers] = useState();

    const getTravelers = useCallback(async () => {
        const res = await GetDataByCrieteries('travelers', {
            'from': fromID.base.value,
            'to': toID.base.value,
            'departure': departure.base.value,
            'arrival': arrival.base.value
        });
        if (res.err && res.err !== "ok") return Notify('fail', "Попутчиков не найдено");
        setTravelers(res)
    }, [fromID, toID, departure, arrival])

    return (
        <SParsels>
            <div className="filters">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>

                <Input type="date" name="departure" base={departure.base} labelText="Выезд:" />
                <Input type="date" name="arrival" base={arrival.base} labelText="Прибытие:" />

                <span className="search_btn" onClick={getTravelers}>
                    <i className="fa fa-search" aria-hidden="true"></i>
                </span>
            </div>

            {
                travelers &&
                <div className="travelers">
                    {travelers?.map(p => <Traveler key={RandomKey()} data={p} />)}
                </div>
            }
        </SParsels>
    )
}