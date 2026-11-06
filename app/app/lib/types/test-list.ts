export interface TestListInput {
  category_id: number;
  value: string;
  label: string;
  sample_amount_liquid?: number;
  sample_amount_solid?: number;
  sample_amount_powder?: number;
  leadtime?: number;
}

export interface TestListItem extends TestListInput {
  xid: string;
  version: number;
  is_current?: boolean;
  category_name: string;
}
