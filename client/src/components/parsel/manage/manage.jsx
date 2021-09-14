import { useCallback, useState } from "react";

import { DbnceCities } from "utils/effects";
import { POSTRequestWithParams } from "utils/api";
import { useInput } from "utils/form";
import { DateFromMilliseconds } from "utils/content";
import { Notify } from "components/app-notification/notification";
import Input from "components/form-input/input";
import SubmitBtn from "components/submit-btn/submit";
import ClipPlash from "components/clips/clips";
import ClippedFiles from "components/clipped-files-plash/plash";

import styled from "styled-components";

const SParsel = styled.form`
    padding: 1rem;
    margin: 1rem;

    & .transit_points,
    & .price_weigth,
    & .title_expire,
    & .contactNumber {
        display: flex;
        align-items: center;
        justify-content: space-between;

        & > * {
            flex-basis: 45%;
        }
    }

    & .photos {
        position: relative;
        margin-bottom: 10rem;
        padding: 1rem;
        display: flex;
        align-items: center;

        &.clipped {
            flex-direction: column;
            align-items: unset;
        }
    }

    @media screen and (max-width: 600px) {
        & .transit_points,
        & .price_weigth,
        & .title_expire,
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

const removeFile = async (id) => {
    const res = await POSTRequestWithParams("/r/image", { 'id': id });
    if (res?.err !== "ok") return Notify('fail', "Фото не удалилось, попробуйте позднее, или сообщите в службу поддрежки")
    return Notify('success', "Фото удалено")
}

export default function ManageParsel({ type = "create", cb, failText, successText, data }) {
    const weight = useInput(data?.weight);
    const price = useInput(data?.price);
    const title = useInput(data?.title);
    const expire = useInput(DateFromMilliseconds(data?.expireDatetime));
    const contactNumber = useInput(data?.contactNumber);
    const from = useInput(data?.from);
    const to = useInput(data?.to);
    const fromID = useInput(data?.fromID);
    const toID = useInput(data?.toID);
    const [isHaveWhatsup, setHaveWhatsup] = useState(data?.isHaveWhatsup === 1);

    from.base.onChange = e => onChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => onChangeTransitPoint(to, e, toID.setCertainValue);

    const [photos, setPhotos] = useState();

    const onSubmit = useCallback(async (e) => {
        e.preventDefault();

        const params = {
            'title': title.base.value,
            'from': fromID.base.value,
            'to': toID.base.value,
            'weight': weight.base.value,
            'price': price.base.value,
            'expireDatetime': Date.parse(expire.base.value),
            'contactNumber': contactNumber.base.value,
            'isHaveWhatsup': isHaveWhatsup ? 1 : 0,
        }
        const res = await POSTRequestWithParams("/" + (type === "create" ? "s" : "e") + "/parsel", params);

        if (res?.err !== "ok") return Notify('fail', failText);
        Notify('success', successText);
        cb(params)
    }, [title, fromID, toID, weight, price, expire, contactNumber, isHaveWhatsup, type, cb, failText, successText]);

    return (
        <SParsel onSubmit={onSubmit}>
            <div className="transit_points">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>
            </div>

            <div className="price_weigth">
                <Input type="number" name="weight" base={weight.base} labelText="Вес (в кг)" />
                <Input type="number" name="price" base={price.base} labelText="Цена (в тг)" />
            </div>

            <div className="title_expire">
                <Input type="text" name="title" base={title.base} labelText="Заголовок вашей посылки" />
                <Input type="date" name="expireDatetime" base={expire.base} labelText="Доставить до:" />
            </div>

            <div className="contactNumber">
                <Input type="tel" name="contactNumber" base={contactNumber.base} labelText="Контакты отправителя" />
                <span>
                    <input onChange={() => setHaveWhatsup(!isHaveWhatsup)} type="checkbox" name="isHaveWhatsup" /> Есть WhatsUp?
                </span>
            </div>

            {
                type === "create"
                    ? <div className="photos">
                        <span>Чтобы прикрепить фото, нажмите здесь {"->"}</span>
                        <ClipPlash setFiles={setPhotos} preloadedFiles={photos} />
                    </div>
                    : <div className="photos clipped">
                        <span>Чтобы удалить фото, нажмите крестик на фото </span>
                        <ClippedFiles files={photos} removeFile={removeFile} />
                    </div>
            }

            <SubmitBtn value={type === "create" ? "Опубликовать" : "Изменить"} />
        </SParsel>
    )
}