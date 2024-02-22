import axios, { AxiosResponse } from 'axios';

const apiUrl: string = 'http://localhost:8080/api/v1/applicationStatus';

export const fetchData = (): Promise<any> => {
    return axios.get(apiUrl)
        .then((response: AxiosResponse<any>) => {
            return response.data;
        })
        .catch((error: any) => {
            console.error('Error fetching data:', error);
            throw new Error('Failed to fetch data');
        });
};
