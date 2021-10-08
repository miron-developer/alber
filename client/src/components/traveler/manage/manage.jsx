import { useCallback, useState } from "react";

import { CompareParams, GetValueFromListByIDAndInputValue, OnChangeTransitPoint } from "utils/effects";
import { POSTRequestWithParams } from "utils/api";
import { useInput } from "utils/form";
import { Notify } from "components/app-notification/notification";
import { ClosePopup } from "components/popup/popup";
import Input from "components/form-input/input";
import SubmitBtn from "components/submit-btn/submit";

import styled from "styled-components";

const STravel = styled.form`
    padding: 1rem;
    margin: 1rem;
    min-width: 80vw;

    & > div {
        margin: 1rem;
    }

    & .transit_points {
        display: flex;
        align-items: center;
        justify-content: space-between;

        & > * {
            flex-basis: 45%;
        }
    }

    & .travel_type {
        display: flex;
        align-items: center;
        justify-content: space-evenly;

        .travel_type-item {
            width: 4rem;
            height: 4rem;
            padding: 1rem;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 2rem;
            color: white;
            border: 1px solid;
            border-radius: 50%;
            transition: .5s;
            cursor: pointer;
            box-shadow: var(--boxShadow);

            &.active,
            &:hover{
                background: #192955;
            }
        }
    }

    @media screen and (max-width: 600px) {
        & .transit_points,
        & .travelType_weigth {
            align-items: unset;
            flex-direction: column;
        }
    }
`;

const clearAll = (fields = [], setHaveWhatsUp, setTravelType) => {
    fields.forEach(f => f.resetField());
    setHaveWhatsUp(false);
    setTravelType(0);
}

const defineTravelTypeByID = (travelTypeID) => {
    if (travelTypeID === 1) return "Машина";
    if (travelTypeID === 2) return "Поезд";
    if (travelTypeID === 3) return "Самолет";
    if (travelTypeID === 4) return "Корабль";
    return ""
}

export default function ManageTraveler({ type = "create", cb, failText = "Ошибка", successText = "Успех", data }) {
    const description = useInput(data?.description);
    const contactNumber = useInput(data?.contactNumber);
    const from = useInput(data?.from);
    const to = useInput(data?.to);
    const travelTypeID = useInput(data?.travelTypeID || 0);
    const fromID = useInput(data?.fromID);
    const toID = useInput(data?.toID);
    const [isHaveWhatsUp, setHaveWhatsUp] = useState(data?.isHaveWhatsUp === 1);

    const [travelType, setTravelType] = useState(data?.travelTypeID || 0);

    from.base.onChange = e => OnChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => OnChangeTransitPoint(to, e, toID.setCertainValue);

    const onSubmit = useCallback(async (e) => {
        e.preventDefault();

        const oldParams = {
            'travelTypeID': data?.travelTypeID,
            'travelType': data?.travelType,
            'fromID': data?.fromID,
            'toID': data?.toID,
            'from': data?.from,
            'to': data?.to,
            'description': data?.description,
            'contactNumber': data?.contactNumber,
            'isHaveWhatsUp': data?.isHaveWhatsUp,
        }
        const comparedParams = CompareParams({
            'id': data?.id,
            'travelTypeID': travelTypeID.base.value,
            'travelType': defineTravelTypeByID(parseInt(travelTypeID.base.value)),
            'fromID': GetValueFromListByIDAndInputValue('from-list', from.base.value),
            'toID': GetValueFromListByIDAndInputValue('to-list', to.base.value),
            'from': from.base.value,
            'to': to.base.value,
            'description': description.base.value,
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
            const fields = [description, contactNumber, travelTypeID, from, to, fromID, toID];
            clearAll(fields, setHaveWhatsUp, setTravelType)
        }
    }, [travelTypeID, from, to, fromID, toID, description, contactNumber, isHaveWhatsUp, type, cb, failText, successText, data]);

    return (
        <STravel onSubmit={onSubmit}>
            <div className="transit_points">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>
            </div>

            <div className="description">
                <textarea
                    className="form-control" {...description.base} required
                    id="description" name="description" rows="3" placeholder="Опишите вашу поездку, сколько вы забираете, когда выходите и приходите"
                ></textarea>
            </div>

            <div className="travel_type">
                <div className={`travel_type-item ${travelType === 1 ? "active" : ""}`} onClick={() => travelTypeID.setCertainValue(1) || setTravelType(1)}>
                    <i className="fa fa-car" aria-hidden="true"></i>
                </div>

                <div className={`travel_type-item ${travelType === 2 ? "active" : ""}`} onClick={() => travelTypeID.setCertainValue(2) || setTravelType(2)}>
                    <i className="fa fa-train" aria-hidden="true"></i>
                </div>

                <div className={`travel_type-item ${travelType === 3 ? "active" : ""}`} onClick={() => travelTypeID.setCertainValue(3) || setTravelType(3)}>
                    <i className="fa fa-plane" aria-hidden="true"></i>
                </div>

                <div className={`travel_type-item ${travelType === 4 ? "active" : ""}`} onClick={() => travelTypeID.setCertainValue(4) || setTravelType(4)}>
                    <i className="fa fa-ship" aria-hidden="true"></i>
                </div>
            </div>

            <Input type="tel" name="contactNumber" base={contactNumber.base} labelText="Номер отправителя" />

            <div className="form-check">
                <label htmlFor="wp" className="form-check-label"></label>
                <input id="wp" className="form-check-input" onChange={() => setHaveWhatsUp(!isHaveWhatsUp)} checked={isHaveWhatsUp} type="checkbox" name="isHaveWhatsup" /> Есть WhatsUp?
            </div>
            <SubmitBtn value={type === "create" ? "Опубликовать" : "Изменить"} />
        </STravel>
    )
}