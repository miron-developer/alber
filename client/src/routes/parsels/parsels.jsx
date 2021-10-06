import { useCallback, useEffect } from "react";

import { GetValueFromListByIDAndInputValue, OnChangeTransitPoint, ScrollHandler } from "utils/effects";
import { useInput } from "utils/form";
import { RandomKey, ValidateParselTravelerSearch } from "utils/content";
import { useFromTo } from "utils/hooks";
import Input from "components/form-input/input";
import Parsel from "components/parsel/parsel";

import styled from "styled-components";

const SParsels = styled.section`
    &.parsels {
        margin-bottom: 10rem;
    }
    
    & .filters {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        justify-content: space-around;
        padding: 1rem;
        background: var(--blueColor);

        & > div {
            flex-basis: 20%;
        }

        & .search_btn {
            padding: .5rem 1rem;
            margin: 0 1rem;
            border-radius: 10px;
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
            
            & > div {
                flex-basis: 50%;
            }

            & .search_btn {
                width: 100%;
                text-align: center;
            }
        }
    }
`;

export default function ParselsPage() {
    const from = useInput('');
    const to = useInput('');
    const fromID = useInput('');
    const toID = useInput('');
    const weight = useInput('');
    const price = useInput('');

    from.base.onChange = e => OnChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => OnChangeTransitPoint(to, e, toID.setCertainValue);

    const { datalist, isStopLoad, getPart } = useFromTo([], 5);

    const loadParsels = useCallback((clear = false) => {
        const params = ValidateParselTravelerSearch(
            GetValueFromListByIDAndInputValue("from-list", from.base.value), GetValueFromListByIDAndInputValue("to-list", to.base.value),
            weight.base.value * 1000, price.base.value
        )
        if (!params) return;
        getPart("parsels", params, 'Не удалось загрузить посылки', true, clear === true ? true : false)
    }, [from, to, weight, price, getPart])


    useEffect(() => {
        // set scroll handler
        document.body.onscroll = e => ScrollHandler(e, isStopLoad, false, loadParsels, "parsel");
    }, [isStopLoad, loadParsels])

    return (
        <SParsels className="parsels">
            <div className="filters">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>

                <Input type="number" name="weight" base={weight.base} labelText="Вес (в кг)" />
                <Input type="number" name="price" base={price.base} labelText="Цена (в тг)" />

                <span className="search_btn btn btn-primary" onClick={() => loadParsels(true)}>
                    <i className="fa fa-search" aria-hidden="true"></i>
                </span>
            </div>

            {
                datalist &&
                <div className="parsels">
                    {datalist?.map(p => <Parsel key={RandomKey()} data={p} />)}
                </div>
            }
        </SParsels>
    )
}