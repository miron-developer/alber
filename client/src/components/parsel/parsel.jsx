import { useCallback, useEffect, useState } from "react";

import { GetDataByCrieteries } from "utils/api";
import { RandomKey } from "utils/content";
import { EditItem, PaintItem, RemoveItem, TopItem } from "utils/effects";
import { Notify } from "components/app-notification/notification";

import styled from "styled-components"

const SParsel = styled.div`
    position: relative;
    padding: 1rem;
    margin: 1rem;
    min-height: 30vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    background: ${props => props.color ? props.color : '#ffffff94'};
    border-radius: 10px;
    box-shadow: var(--boxShadow);

    & .info {
        display: flex;
        justify-content: space-between;

        & .general_info {
            display: flex;
            flex-direction: column;
            flex-basis: 70%;

            & span {
                word-break: break-word;
            }

            & .price b {
                text-decoration: underline;
            }
        }

        & .other_info  {
            & .phones > * {
                margin: .5rem;
                padding: .5rem;
                font-size: 4rem;
                border-radius: 10px;
                cursor: pointer;
                transition: var(--transitionApp);

                &:hover {
                    color: var(--onHoverColor);
                    background: var(--blueColor);
                }
            }

            & .expire {
                margin: .5rem;
            }
        }
    }

    & .photos {
        display: flex;
        flex-wrap: wrap;

        & span {
            max-width: 20vw;
            display: block;
            align-self: end;
            margin: 5px;
            border-radius: 5px;
            background: #000000a3;

            & img {
                width: 100%;
            }
        }
    }

    & .manage {
        position: absolute;
        right: 0;
        top: 0;
        padding: .5rem;
        margin: .5rem;    
        background: #2f3a64;
        color: var(--onHoverColor);
        font-size: 1rem;
        border-radius: 5px;

        & .manage-action {
            cursor: pointer;
            margin: .5rem;
            transition: var(--transitionApp);
            text-align: right;
            vertical-align: middle;
            
            &:hover {
                background: var(--darkGreyColor);
            }
        }

        & .manage-actions {
            display: flex;
            flex-direction: column;
        }
    }

    @media screen and (max-width: 600px) {
        & .info {
            flex-direction: column;

            & .phones {
                display: flex;
                align-items: center;
                justify-content: space-evenly;
            }
        }

        & .manage {
            font-size: 1.5rem;
        }
    }
`;

// Paint
// <span className="manage-action" onClick={() => PaintItem(data.id, "parsel", newData => changeItem(data.id, Object.assign({}, data, newData)))}>
//     <i className="fa fa-paint-brush" aria-hidden="true">Покрасить</i>
// </span>

export default function Parsel({ data, isMy = false, changeItem, removeItem, reloadItem }) {
    const [photos, setPhotos] = useState();
    const [isOpened, setOpened] = useState(false);

    const getPhotos = useCallback(async () => {
        const res = await GetDataByCrieteries("images", { "id": data.id });
        if (res?.err === "н/д") return setPhotos(null);
        if (res.err && res.err !== "ok") return Notify('fail', "Не удалось загрузить прикрепленные фото");
        return setPhotos(res?.data || res);
    }, [data]);

    useEffect(() => {
        if (photos === undefined) return getPhotos()
    }, [getPhotos, photos])

    return (
        <SParsel className="parsel" color={data.color}>
            <div className="info">
                <div className="general_info">
                    <span>Имя: {data.nickname}</span>
                    <span>{data.from}-{data.to}</span>
                    <span>{data.title}</span>
                    <span>Вес: {data.weight / 1000} кг</span>
                    <span className="price">Цена: <b> {data.price} </b> тг</span>
                    <span>Описание: {data.description}</span>
                </div>

                <div className="other_info">
                    <div className="phones">
                        {
                            data.isHaveWhatsUp === 1 &&
                            <a target="_blank" rel="noreferrer" href={`https://api.whatsapp.com/send?phone=${data.contactNumber}&text="Добрый день, пишу из приложения Al-Ber насчет вашей посылки: ${data.title}"`}>
                                <i className="fa fa-whatsapp" aria-hidden="true"></i>
                            </a>
                        }

                        <span onClick={() => window.open("tel:" + data.contactNumber)}>
                            <i className="fa fa-phone" aria-hidden="true"></i>
                        </span>
                    </div>
                </div>
            </div>

            {
                photos
                    ? <div className="photos">
                        {photos?.map(({ src }) => <span key={RandomKey()}><img src={src} alt="" /></span>)}
                    </div>
                    : <div className="photos"> Нет фото</div>
            }

            {
                isMy &&
                <div className="manage">
                    <div className="manage-action" onClick={() => setOpened(!isOpened)}>
                        <span>Действия</span>
                    </div>

                    {
                        isOpened &&
                        <div className="manage-actions">
                            <span className="manage-action"
                                onClick={() => EditItem("parsel", { ...data, 'photos': photos }, newData => changeItem(data.id, newData), reloadItem)}
                            >
                                <i className="fa fa-pencil" aria-hidden="true">Редактировать</i>
                            </span>
                            <span className="manage-action" onClick={() => RemoveItem(data.id, "parsel", () => removeItem(data.id))}>
                                <i className="fa fa-trash" aria-hidden="true">Удалить</i>
                            </span>

                            <span className="manage-action" onClick={() => TopItem(data.id, "parsel", newData => changeItem(data.id, Object.assign({}, data, newData)))}>
                                <i className="fa fa-level-up" aria-hidden="true">Поднять</i>
                            </span>
                        </div>
                    }
                </div>
            }
        </SParsel>
    )
}