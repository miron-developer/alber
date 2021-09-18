import { useCallback } from "react";

import { GetValueFromListByIDAndInputValue, OnChangeTransitPoint, ScrollHandler } from "utils/effects";
import { useInput } from "utils/form";
import { RandomKey, ValidateParselTravelerSearch } from "utils/content";
import { useFromTo } from "utils/hooks";
import Input from "components/form-input/input";
import Traveler from "components/traveler/traveler";

import styled from "styled-components";

const STravelers = styled.section`
    & .filters {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        justify-content: space-evenly;
        padding: 1rem;
        background: var(--blueColor);

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

    const { datalist, isStopLoad, getPart } = useFromTo([], 5)

    const loadTravelers = useCallback((clear = false) => {
        const params = ValidateParselTravelerSearch(
            GetValueFromListByIDAndInputValue("from-list", from.base.value), GetValueFromListByIDAndInputValue("to-list", to.base.value),
            Date.parse(departure.base.value), Date.parse(arrival.base.value)
        )
        if (!params) return;
        getPart("travelers", params, 'Не удалось загрузить попутчиков', true, clear === true ? true : false)
    }, [from, to, departure, arrival, getPart])

    // set scroll handler
    document.body.onscroll = e => ScrollHandler(e, isStopLoad, false, loadTravelers, "traveler");

    return (
        <STravelers>
            <div className="filters">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>

                <Input type="date" name="departure" base={departure.base} labelText="Выезд:" />
                <Input type="date" name="arrival" base={arrival.base} labelText="Прибытие:" />

                <span className="search_btn" onClick={() => loadTravelers(true)}>
                    <i className="fa fa-search" aria-hidden="true"></i>
                </span>
            </div>

            {
                datalist &&
                <div className="travelers">
                    {datalist?.map(p => <Traveler key={RandomKey()} data={p} />)}
                </div>
            }
        </STravelers>
    )
}