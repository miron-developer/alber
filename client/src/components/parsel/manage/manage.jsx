import { useCallback, useState } from "react";

import { CompareParams, GetValueFromListByIDAndInputValue, OnChangeTransitPoint } from "utils/effects";
import { POSTRequestWithParams } from "utils/api";
import { useInput } from "utils/form";
import { UploadFile } from "utils/file";
import { DateFromMilliseconds } from "utils/content";
import { Notify } from "components/app-notification/notification";
import { ClosePopup } from "components/popup/popup";
import Input from "components/form-input/input";
import SubmitBtn from "components/submit-btn/submit";
import ClipPlash from "components/clips/clips";
import ClippedFiles from "components/clipped-files-plash/plash";

import styled from "styled-components";

const SParsel = styled.form`
    padding: 1rem;
    margin: 1rem;
    min-width: 80vw;

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
            margin-bottom: 0;
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

const removeFile = async (id, src, removePhoto) => {
    const res = await POSTRequestWithParams("/r/image", { 'id': id, 'src': src });
    if (res?.err !== "ok") return Notify('fail', "Фото не удалилось, попробуйте позднее, или сообщите в службу поддрежки")
    removePhoto(id);
    return Notify('success', "Фото удалено")
}

const clearAll = (fields = [], setHaveWhatsUp, setPreloaded) => {
    fields.forEach(f => f.resetField());
    setHaveWhatsUp(false);
    setPreloaded([]);
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
    const [isHaveWhatsUp, setHaveWhatsUp] = useState(data?.isHaveWhatsUp === 1);

    from.base.onChange = e => OnChangeTransitPoint(from, e, fromID.setCertainValue);
    to.base.onChange = e => OnChangeTransitPoint(to, e, toID.setCertainValue);

    const [photos, setPhotos] = useState(data?.photos);
    const [preloadedFiles, setPreloaded] = useState([]);

    const removePhoto = id => setPhotos(photos.filter(ph => ph.id !== id))

    const onSubmit = useCallback(async (e) => {
        e.preventDefault();

        const oldParams = {
            'title': data?.title,
            'fromID': data?.fromID,
            'toID': data?.toID,
            'from': data?.from,
            'to': data?.to,
            'weight': data?.weight,
            'price': data?.price,
            'expireDatetime': data?.expireDatetime,
            'contactNumber': data?.contactNumber,
            'isHaveWhatsUp': data?.isHaveWhatsUp,
        }
        const comparedParams = CompareParams({
            'id': data?.id,
            'title': title.base.value,
            'fromID': GetValueFromListByIDAndInputValue('from-list', from.base.value),
            'toID': GetValueFromListByIDAndInputValue('to-list', to.base.value),
            'from': from.base.value,
            'to': to.base.value,
            'weight': weight.base.value,
            'price': price.base.value,
            'expireDatetime': Date.parse(expire.base.value),
            'contactNumber': contactNumber.base.value,
            'isHaveWhatsUp': isHaveWhatsUp ? 1 : 0,
        }, oldParams);

        // bcs we have id on new so <= 1
        if (Object.values(comparedParams).length <= 1 && preloadedFiles.length === 0) return Notify('info', 'Нет изменений');

        // send to edit
        const res = await POSTRequestWithParams("/" + (type === "create" ? "s" : "e") + "/parsel", comparedParams);
        if (res?.err !== "ok") return Notify('fail', failText + ":" + res?.err);
        Notify('success', successText);

        // upload images
        preloadedFiles.forEach(ph => UploadFile(ph.type, ph.file, "parsel", type === "create" ? res?.data : data?.id))

        // do callback if edit
        if (cb) {
            // finally params will be
            cb(Object.assign(oldParams, comparedParams));
            ClosePopup();
        } else {
            // or clear all if create
            const fields = [weight, price, title, expire, contactNumber, from, to, fromID, toID];
            clearAll(fields, setHaveWhatsUp, setPreloaded)
        }
    }, [title, from, to, fromID, toID, weight, price, expire, contactNumber, isHaveWhatsUp, preloadedFiles, type, cb, failText, successText, data]);

    return (
        <SParsel onSubmit={onSubmit}>
            <div className="transit_points">
                <Input id="from" type="text" name="from" list="from-list" base={from.base} labelText="Откуда" />
                <datalist id="from-list"></datalist>

                <Input id="to" type="text" name="to" list="to-list" base={to.base} labelText="Куда" />
                <datalist id="to-list"></datalist>
            </div>

            <div className="price_weigth">
                <Input type="number" name="weight" base={weight.base} labelText="Вес (в г)" />
                <Input type="number" name="price" base={price.base} labelText="Цена (в тг)" />
            </div>

            <div className="title_expire">
                <Input type="text" name="title" base={title.base} labelText="Заголовок вашей посылки" />
                <Input type="datetime-local" name="expireDatetime" base={expire.base} labelText="Доставить до:" />
            </div>

            <div className="contactNumber">
                <Input type="tel" name="contactNumber" base={contactNumber.base} labelText="Контакты отправителя" />
                <span>
                    <input onChange={() => setHaveWhatsUp(!isHaveWhatsUp)} checked={isHaveWhatsUp} type="checkbox" name="isHaveWhatsup" /> Есть WhatsUp?
                </span>
            </div>

            {
                type !== "create" && photos && photos.length > 0 &&
                <div className="photos clipped">
                    <span>Чтобы удалить фото, нажмите крестик на фото </span>
                    <ClippedFiles files={photos} removeFile={(id, src) => removeFile(id, src, removePhoto)} />
                </div>
            }

            <div className="photos">
                <span>Чтобы прикрепить фото, нажмите здесь {"->"}</span>
                <ClipPlash setFiles={setPreloaded} preloadedFiles={preloadedFiles} />
            </div>

            <SubmitBtn value={type === "create" ? "Опубликовать" : "Изменить"} />
        </SParsel>
    )
}