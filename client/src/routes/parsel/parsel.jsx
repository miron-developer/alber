import { USER } from "constants/constants"
import ManageParsel from "components/parsel/manage/manage";

import styled from "styled-components";

const SParsel = styled.section`
    padding: 1rem;
    margin: 1rem;

    & .parsel_create_tip {
        padding: 1rem;
        margin: 1rem;
        border-radius: 10px;
        text-align: center;
        background: #f5f000;
        box-shadow: var(--boxShadow);
    }
`;

const SParselGuest = styled.section`
    padding: 1rem;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`;

export default function ParselPage() {
    if (USER.status !== "online" || USER.guest) return <SParselGuest>Войдите, чтобы создать посылки</SParselGuest>
    return (
        <SParsel>
            <div className="parsel_create_tip">
                Заполните данные, и ваша посылка попадет в ленту посылок, оттуда люди могут увидеть и забрать Вашу посылку
            </div>

            <ManageParsel failText="Не удалось создать посылку" successText="Создана посылка" />
        </SParsel>
    )
}