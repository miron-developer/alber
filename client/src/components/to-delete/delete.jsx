import { POSTRequestWithParams } from "utils/api";
import { ClosePopup } from "components/popup/popup";
import { Notify } from "components/app-notification/notification";

import styled from "styled-components";

const SToUp = styled.div`
    padding: 1rem;
    margin: 1rem;

    & .price {
        color: red;
        font-size: 1.3rem;
        background: yellow;
    }

    & .answer {
        display: flex;
        align-items: center;
        justify-content: space-evenly;

        & span {
            padding: .5rem 1rem;
            margin: 1rem;
            color: var(--onHoverColor);
            background: var(--blueColor);
            border-radius: 10px;
            cursor: pointer;
            transition: var(--transtionApp);

            &:nth-child(2) {
                background: red;
            }

            &:hover{
                background: var(--onHoverBG);
            }
        }
    }
`;

const toDelete = async (id, type, cb) => {
    const res = await POSTRequestWithParams("/r/" + (type === "parsel" ? "parsel" : 'travel'), { 'id': id })
    if (res.err && res.err !== "ok") return Notify('fail', 'Не удалено');
    cb()
    ClosePopup();
}

/**
 * 
 * @param type if cost will be relative by type: parsel or travel 
 * @param cb callback after click to delete
 * @param id parsel/travel id 
 */
export default function ToDelete({ cb, type, id }) {
    return (
        <SToUp>
            <h2>Удалить Ваше объявление?</h2>
            <span>Удаляется без возвратно</span>

            <div className="answer">
                <span onClick={() => toDelete(id, type, cb)}>Да</span>
                <span onClick={ClosePopup}>Нет</span>
            </div>
        </SToUp>
    )
}