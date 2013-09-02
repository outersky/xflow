package controllers

type AjaxResult struct {
    Success bool
    Validation interface{}
    Meta interface{}
    Result interface{}
    TotalCount int
}

func Error(err interface{}) *AjaxResult {
    result := &AjaxResult{}
    result.Validation = err
    result.Success = false
    return result
}

func PageList(data []interface{}, total int) *AjaxResult {
    result := List(data)
    result.TotalCount = total
    return result
}

func List(data []interface{}) *AjaxResult {
    result := &AjaxResult{}
    result.Result = data
    result.Success = true
    result.TotalCount = len(data)
    return result
}

func Single(data interface{}) *AjaxResult {
    result := &AjaxResult{}
    result.Result = data
    result.Success = true
    result.TotalCount = 1
    return result
}
