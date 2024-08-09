package modules

//
//func Test_AiSyncSetDataToBody(t *testing.T) {
//	param := &models.InputParam{
//		ApiParam: &models.ApiParam{
//			HttpMethod: "POST",
//			Path:       "/api/v10/getElementStyle",
//			Body: map[string]interface{}{
//				"lso":  123,
//				"mosi": 123,
//				"object": map[string]interface{}{
//					"sad":    123,
//					"opwwew": map[string]interface{}{},
//				},
//			},
//		},
//	}
//
//	syncInfo := &ApiSyncBuilder{ApiSync: ApiSync{
//		HttpMethod: "POST",
//		PageSize:   20,
//		PageFiledPosition: &PageFiledPosition{
//			Type:          PageFiledType(1),
//			NumPosition:   "object.opwwew.current",
//			CountPosition: "object.opwwew.pageSize",
//		},
//	},
//	}
//
//	param.ApiParam.Body = utils.SetFieldToMapData(
//		reflect.ValueOf(param.GetApiParam().GetBody()),
//		reflect.ValueOf(0),
//		syncInfo.GetPageFiledPosition().GetNumSlice(), 0)
//
//	fmt.Println(param.ApiParam.Body)
//
//}
