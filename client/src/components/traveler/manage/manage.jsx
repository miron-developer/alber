import { useCallback, useEffect, useState } from "react";

import { CompareParams, DbnceCities } from "utils/effects";
import { GetDataByCrieteries, POSTRequestWithParams } from "utils/api";
import { useInput } from "utils/form";
import { DateFromMilliseconds } from "utils/content";
import { Notify } from "components/app-notification/notification";
import Input from "components/form-input/input";
import SubmitBtn from "components/submit-btn/submit";

import styled from "styled-components";
import Select from "components/form-select/select";

const STravel = styled.form`
    padding: 1rem;
    margin: 1rem;
    min-width: 80vw;

    & .transit_points,
    & .travelType_weigth,
    & .departure_arrival,
    & .contactNumber {
        display: flex;
        align-items: center;
        justify-content: space-between;

        & > * {
            flex-basis: 45%;
        }
    }

    @media screen and (max-width: 600px) {
        & .transit_points,
        & .travelType_weigth,
        & .departure_arrival,
        & .contactNumber {
            align-items: unset;
            flex-direction: column;
        }
    }
`;

const onChangeTransitPoint = async (point, e, setID) => {
    point.setCertainValue(e.target.value);
    DbnceCities(e);
    const dt = Array.from(document.getElementById(e.target.list.id).childNodes)
    if (dt.length === 0) return;
    const op = dt.find(option => option?.value?.includes(e.target.value));
    if (op) setID(op.textContent);
}

const onChangeTravelType = (e, setID) => setID(e.target.value)

export default function ManageTraveler({ type = "create", cb, failText, successText, data }) {
    const weight = useInput(data?.weight);
    const departure = useInput(DateFromMilliseconds(data?.departureDatetime));
    const arrival = useInput(DateFromMilliseconds(data?.arrivalDatetime));
    const contactNumber = useInput(data?.contactNumber);
    const from = useInput(data?.from);
    const to = useInput(data?.to);
    const travelTypeID = useInput(data?.travelTypeID);
    const fromID = useInput(data?.fromID);
    const toID = useInput(data?.toID);
    const [isHaveWhatsup, setHaveWhatsup] = useState(data?.isHaveWhatsup === 1);
    const [travelTypes, setTravelTypes] = useState();

    from.base.onChange = e => onChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => onChangeTransitPoint(to, e, toID.setCertainValue);

    const getTravelTypes = useCallback(async () => {
        const res = await GetDataByCrieteries('travelTypes');
        if (res.err && res?.err !== "ok") return setTravelTypes(null);
        return setTravelTypes(res)
    }, [])

    const onSubmit = useCallback(async (e) => {
        e.preventDefault();

        const comparedParams = CompareParams({
            'travelType': travelTypeID.base.value,
            'from': fromID.base.value,
            'to': toID.base.value,
            'weight': weight.base.value,
            'departureDatetime': Date.parse(departure.base.value),
            'arrivalDatetime': Date.parse(arrival.base.value),
            'contactNumber': contactNumber.base.value,
            'isHaveWhatsup': isHaveWhatsup ? 1 : 0,
        }, {
            'travelType': data?.travelTypeID,
            'from': data?.fromID,
            'to': data?.toID,
            'weight': data?.weight,
            'departureDatetime': data?.departureDatetime,
            'arrivalDatetime': data?.arrivalDatetime,
            'contactNumber': data?.contactNumber,
            'isHaveWhatsup': data?.isHaveWhatsup,
        });
        if (Object.values(comparedParams).length === 0) return Notify('info', 'Нет изменений');
        const res = await POSTRequestWithParams("/" + (type === "create" ? "s" : "e") + "/travel", comparedParams);

        if (res?.err !== "ok") return Notify('fail', failText);
        Notify('success', successText);
        cb(comparedParams)
    }, [travelTypeID, fromID, toID, weight, departure, arrival, contactNumber, isHaveWhatsup, type, cb, failText, successText, data]);

    useEffect(() => {
        if (travelTypes === undefined) return getTravelTypes()
    }, [getTravelTypes, travelTypes])

    console.log('trvtype', travelTypeID);

    return (
        <STravel onSubmit={onSubmit}>
            <div className="transit_points">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>
            </div>

            <div className="travelType_weigth">
                <Input type="number" name="weight" base={weight.base} labelText="Вес (в кг)" />
                <Select name="travelType" text="Тип транспорта" options={{
                    data: travelTypes,
                    value: "id",
                    selected: travelTypeID.base.value,
                    makeText: ({ name }) => name
                }} onChange={e=>onChangeTravelType(e, travelTypeID.setCertainValue)} />
            </div>

            <div className="departure_arrival">
                <Input type="date" name="departureDatetime" base={departure.base} labelText="Выезд" />
                <Input type="date" name="arrivalDatetime" base={arrival.base} labelText="Прибытие" />
            </div>

            <div className="contactNumber">
                <Input type="tel" name="contactNumber" base={contactNumber.base} labelText="Контакты отправителя" />
                <span>
                    <input onChange={() => setHaveWhatsup(!isHaveWhatsup)} type="checkbox" name="isHaveWhatsup" /> Есть WhatsUp?
                </span>
            </div>

            <SubmitBtn value={type === "create" ? "Опубликовать" : "Изменить"} />
        </STravel>
    )
}