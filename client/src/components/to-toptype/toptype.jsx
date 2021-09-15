
import { useCallback, useEffect, useState } from "react";

import { GetDataByCrieteries, POSTRequestWithParams } from "utils/api";
import { useInput } from "utils/form";
import { RandomKey } from "utils/content";
import { Notify } from "components/app-notification/notification";
import Input from "components/form-input/input";
import SubmitBtn from "components/submit-btn/submit";

import styled from "styled-components";

const SToTopType = styled.div`
    padding: 1rem;
    margin: 1rem;
`;

const SOneItem = styled.div`
    padding: 1rem;
    margin: 1rem;
    display: flex;
    flex-direction: column;
    background: ${props => props.color ? props.color : 'var(--grayColor)'};
    border-radius: 5px;
    box-shadow: var(--boxShadow);
    cursor: pointer;
    transition: var(--transitionApp);

    &:hover {
        color: var(--onHoverColor);
        background: var(--onHoverBG);
    }

    & b {
        display: block;
        padding: 5px;
        margin: 5px;
        color: var(--onHoverColor);
        background: #c30000;
        border-radius: 5px;
        box-shadow: var(--boxShadow);
    }
`;

const toTopType = async (id, type, code, topID, cb) => {
    const res = await POSTRequestWithParams("/e/toptype", { 'id': id, 'type': type, 'code': code, 'topID': topID })
    if (res.err && res.err !== "ok") return Notify('fail', 'Не удалось рекламировать');
    cb()
}

const GOnePrice = ({ id, cost, name, color, setPayed, setTopTypeID }) => {
    const onClick = () => {
        setPayed(true);
        setTopTypeID(id);
    }

    return (
        <SOneItem color={color} onClick={onClick}>
            <span>Ваше обьявление будет рекламировано <b> {name} </b></span>
            <span>Стоимость: <b> {cost} тг </b></span>
        </SOneItem>
    )
}

/**
 * 
 * @param type if cost will be relative by type: parsel or travel 
 * @param cb callback after click to up
 * @param id parsel/travel id 
 */
export default function ToTopType({ cb, type, id }) {
    const [prices, setPrices] = useState();
    const [isPayed, setPayed] = useState();
    const [topTypeID, setTopTypeID] = useState();
    const code = useInput('');

    const getPrices = useCallback(async () => {
        const res = await GetDataByCrieteries('toptypes');
        if (res.err && res.err !== "ok") {
            setPrices(null);
            return Notify('fail', "Ошибка. Попробуйте позднее");
        }
        setPrices(res)
    }, [])

    useEffect(() => {
        if (prices === undefined) return getPrices()
    }, [getPrices, prices])


    if (!prices) return <div>Ошибка. Попробуйте позднее</div>
    return (
        <SToTopType>
            {
                isPayed
                    ? <div>
                        <Input index="2" id="code" type="text" name="code" base={code.base} labelText="Введите 8-значный код:"
                            minLength="8" maxLength="8" placeholder="Mfa7sd45"
                        />

                        <SubmitBtn value="Рекламировать!" onClick={() => toTopType(id, type, code, topTypeID, cb)} />
                    </div>
                    : <>
                        <h2>Выберите на какой промежуток Вы хотите рекламировать</h2>

                        <div className="prices">
                            {prices?.map(p => <GOnePrice key={RandomKey()} {...p} setPayed={setPayed} setTopTypeID={setTopTypeID} />)}
                        </div>
                    </>
            }
        </SToTopType>
    )
}