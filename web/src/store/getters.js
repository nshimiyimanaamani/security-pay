const getters = {
  getEndpoint: state => state.endpoint,
  getSectorArray: state => state.sector,
  getCellsArray: state => state.cells_array,
  getActiveCell: state => state.active_cell,
  getActiveVillage: state => state.active_village,
  getVillageArray: state => state.village_array,
  getActiveSector: state => state.active_sector,
  userDetails: state => state.user,
  getMonths: state => state.months,
  location: state => {
    return {
      province: state.province,
      district: state.district,
      sector: state.active_sector,
      cell: state.active_cell,
      village: state.active_village
    };
  },
  appLoading: state => state.appLoading
};

export default getters;
