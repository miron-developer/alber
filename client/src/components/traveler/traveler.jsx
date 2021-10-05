import { useState } from "react";

import { EditItem, PaintItem, RemoveItem, TopItem } from "utils/effects";

import styled from "styled-components";

const STraveler = styled.div`
    position: relative;
    padding: 1rem;
    margin: 1rem;
    min-height: 30vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    background: ${props => props.color ? props.color : '#ffffff94'};
    border-radius: 10px;

    & .info {
        display: flex;
        justify-content: space-between;

        & .general_info {
            display: flex;

            & span {
                word-break: break-word;
            }

            & div {
                display: flex;
                flex-direction: column;
                margin: 1rem;
            }

            & .weigth b {
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

            & .general_info {
                flex-direction: column;
            }

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

// paint
// <span className="manage-action" onClick={() => PaintItem(data.id, "traveler", newData => changeItem(data.id, Object.assign({}, data, newData)))}>
//     <i className="fa fa-paint-brush" aria-hidden="true">Покрасить</i>
// </span>

export default function Traveler({ data, isMy = false, changeItem, removeItem }) {
    const [isOpened, setOpened] = useState(false);

    return (
        <STraveler className="traveler" color={data.color}>
            <div className="info">
                <div className="general_info">
                    <div className="common">
                        <span>Имя: {data.nickname}</span>
                        <span>{data.from}-{data.to}</span>
                        <span>Тип транспорта: {data.travelType}</span>
                        <span>Описание: {data.description}</span>
                    </div>
                </div>

                <div className="other_info">
                    <div className="phones">
                        {
                            data.isHaveWhatsUp === 1 &&
                            <a target="_blank" rel="noreferrer" href={`https://api.whatsapp.com/send?phone=${data.contactNumber}&text="Добрый день, пишу из приложения Al-Ber насчет передачи посылки"`}>
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
                isMy &&
                <div className="manage">
                    <div className="manage-action" onClick={() => setOpened(!isOpened)}>
                        <span>Действия</span>
                    </div>

                    {
                        isOpened &&
                        <div className="manage-actions">
                            <span className="manage-action"
                                onClick={
                                    () =>
                                        EditItem(
                                            "traveler",
                                            data,
                                            newData => changeItem(data.id, newData)
                                        )
                                }
                            >
                                <i className="fa fa-pencil" aria-hidden="true">Редактировать</i>
                            </span>
                            <span className="manage-action" onClick={() => RemoveItem(data.id, "traveler", () => removeItem(data.id))}>
                                <i className="fa fa-trash" aria-hidden="true">Удалить</i>
                            </span>

                            <span className="manage-action" onClick={() => TopItem(data.id, "traveler", newData => changeItem(data.id, Object.assign({}, data, newData)))}>
                                <i className="fa fa-level-up" aria-hidden="true">Поднять</i>
                            </span>
                        </div>
                    }
                </div>
            }
        </STraveler>
    )
}