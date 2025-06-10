import axios from "axios";

export type GetAutocompleteLogParams = {
  kind: "label"|"identity"
  key: string
  value: string
};

export function getAutocompleteLog(params: GetAutocompleteLogParams): Promise<string[]> {
  return axios.get<string[]>("/api/autocomplete/log", { params: params}).then((resp) => {
    return resp.data;
  });
}

