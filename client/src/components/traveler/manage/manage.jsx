import { useCallback, useEffect, useState } from "react";

import { CompareParams, GetValueFromListByIDAndInputValue, OnChangeTransitPoint } from "utils/effects";
import { GetDataByCrieteries, POSTRequestWithParams } from "utils/api";
import { useInput } from "utils/form";
import { DateFromMilliseconds } from "utils/content";
import { Notify } from "components/app-notification/notification";
import { ClosePopup } from "components/popup/popup";
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

const onChangeTravelType = (e, setID, setTravel) => {
    console.log(e);
    setID(e.target.value);
    const op = Array.from(e.target.options).find(op => op.selected)
    if (op) setTravel(op.textContent);
}

const clearAll = (fields = [], setHaveWhatsUp) => {
    fields.forEach(f => f.resetField());
    setHaveWhatsUp(false);
}

export default function ManageTraveler({ type = "create", cb, failText, successText, data }) {
    const weight = useInput(data?.weight);
    const departure = useInput(DateFromMilliseconds(data?.departureDatetime));
    const arrival = useInput(DateFromMilliseconds(data?.arrivalDatetime));
    const contactNumber = useInput(data?.contactNumber);
    const from = useInput(data?.from);
    const to = useInput(data?.to);
    const travelType = useInput(data?.travelType);
    const travelTypeID = useInput(data?.travelTypeID || 1);
    const fromID = useInput(data?.fromID);
    const toID = useInput(data?.toID);
    const [isHaveWhatsUp, setHaveWhatsUp] = useState(data?.isHaveWhatsup === 1);
    const [travelTypes, setTravelTypes] = useState();

    from.base.onChange = e => OnChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => OnChangeTransitPoint(to, e, toID.setCertainValue);

    const getTravelTypes = useCallback(async () => {
        const res = await GetDataByCrieteries('travelTypes');
        if (res.err && res?.err !== "ok") return setTravelTypes(null);
        return setTravelTypes(res)
    }, [])

    const onSubmit = useCallback(async (e) => {
        e.preventDefault();

        const oldParams = {
            'travelType': data?.travelType,
            'travelTypeID': data?.travelTypeID,
            'fromID': data?.fromID,
            'toID': data?.toID,
            'from': data?.from,
            'to': data?.to,
            'weight': data?.weight,
            'departureDatetime': data?.departureDatetime,
            'arrivalDatetime': data?.arrivalDatetime,
            'contactNumber': data?.contactNumber,
            'isHaveWhatsUp': data?.isHaveWhatsUp,
        }
        const comparedParams = CompareParams({
            'id': data?.id,
            'travelType': travelType.base.value,
            'travelTypeID': travelTypeID.base.value,
            'fromID': GetValueFromListByIDAndInputValue('from-list', from.base.value),
            'toID': GetValueFromListByIDAndInputValue('to-list', to.base.value),
            'from': from.base.value,
            'to': to.base.value,
            'weight': weight.base.value,
            'departureDatetime': Date.parse(departure.base.value),
            'arrivalDatetime': Date.parse(arrival.base.value),
            'contactNumber': contactNumber.base.value,
            'isHaveWhatsUp': isHaveWhatsUp ? 1 : 0,
        }, oldParams);

        // bcs we have id on new, so <= 1
        if (Object.values(comparedParams).length <= 1) return Notify('info', 'Нет изменений');

        // send
        const res = await POSTRequestWithParams("/" + (type === "create" ? "s" : "e") + "/travel", comparedParams);
        if (res?.err !== "ok") return Notify('fail', failText + ":" + res?.err);
        Notify('success', successText);

        // do callback if edit
        if (cb) {
            // finally params will be
            cb(Object.assign(oldParams, comparedParams));
            ClosePopup()
        } else {
            // or clear all if create
            const fields = [weight, departure, arrival, travelTypeID, travelType, contactNumber, from, to, fromID, toID];
            clearAll(fields, setHaveWhatsUp)
        }
    }, [travelTypeID, travelType, from, to, fromID, toID, weight, departure, arrival, contactNumber, isHaveWhatsUp, type, cb, failText, successText, data]);

    useEffect(() => {
        if (travelTypes === undefined) return getTravelTypes()
    }, [getTravelTypes, travelTypes])

    return (
        <STravel onSubmit={onSubmit}>
            <div className="transit_points">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>
            </div>

            <div className="travelType_weigth">
                <Input type="number" name="weight" base={weight.base} labelText="Заберу до (в г)" />
                <Select name="travelType" text="Тип транспорта" options={{
                    data: travelTypes,
                    value: "id",
                    selected: travelTypeID.base.value,
                    makeText: ({ name }) => name
                }} onChange={e => onChangeTravelType(e, travelTypeID.setCertainValue, travelType.setCertainValue)} />
            </div>

            <div className="departure_arrival">
                <Input type="datetime-local" name="departureDatetime" base={departure.base} labelText="Выезд" />
                <Input type="datetime-local" name="arrivalDatetime" base={arrival.base} labelText="Прибытие" />
            </div>

            <div className="contactNumber">
                <Input type="tel" name="contactNumber" base={contactNumber.base} labelText="Контакты отправителя" />
                <span>
                    <input onChange={() => setHaveWhatsUp(!isHaveWhatsUp)} type="checkbox" name="isHaveWhatsUp" /> Есть WhatsUp?
                </span>
            </div>

            <SubmitBtn value={type === "create" ? "Опубликовать" : "Изменить"} />
        </STravel>
    )
}