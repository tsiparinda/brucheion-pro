package repo

import "brucheion/models"

// func (repo *SqlRepository) GetProduct(id int) (p models.Product) {
//     row := repo.Commands.GetProduct.QueryRowContext(repo.Context, id)
//     if row.Err() == nil {
//         var err error
//         if p, err = scanProduct(row); err != nil {
//             repo.Logger.Panicf("Cannot scan data: %v", err.Error())    
//         }
//     } else {
//         repo.Logger.Panicf("Cannot exec GetProduct command: %v", row.Err().Error())
//     }
//     return
// }

func (repo *SqlRepository) GetBoltCatalog() (results []models.BoltCatalog) {
    rows, err := repo.Commands.GetBoltCatalog.QueryContext(repo.Context)
    if err == nil {
        if results, err = scanBoltCatalog(rows); err != nil {
            repo.Logger.Panicf("Cannot scan data: %v", err.Error())                    
            return
        }
    } else {
        repo.Logger.Panicf("Cannot exec GetProducts command: %v", err)   
    }
    return
}

