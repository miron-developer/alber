import { USER } from "constants/constants"
import ManageTraveler from "components/traveler/manage/manage";

import styled from "styled-components";

const STravel = styled.section`
    padding: 1rem;
    margin: 1rem;

    & .travel_create_tip {
        padding: 1rem;
        margin: 1rem;
        border-radius: 10px;
        text-align: center;
        background: #f5f000;
        box-shadow: var(--boxShadow);
    }
`;

const STravelGuest = styled.section`
    padding: 1rem;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`;

export default function TravelerPage() {
    if (USER.status !== "online" || USER.guest) return <STravelGuest>Войдите, чтобы стать попутчиком</STravelGuest>
    return (
        <STravel>
            <div className="travel_create_tip">
                Заполните данные, и Ваша заявка попадет в ленту попутчиков, оттуда Вас заметят люди, которые хотят отправить посылку, а Вы заработаете
            </div>

            <ManageTraveler failText="Не удалось создать объявление попутчика" successText="Объявление попутчика создан" />
        </STravel>
    )
}