import { TestBed, async } from '@angular/core/testing';
import { OnosComponent } from './onos.component';

describe('OnosComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        OnosComponent
      ],
    }).compileComponents();
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(OnosComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have as title 'sd-ran-gui'`, () => {
    const fixture = TestBed.createComponent(OnosComponent);
    const app = fixture.componentInstance;
    expect(app.title).toEqual('sd-ran-gui');
  });

  it('should render title', () => {
    const fixture = TestBed.createComponent(OnosComponent);
    fixture.detectChanges();
    const compiled = fixture.nativeElement;
    expect(compiled.querySelector('.content span').textContent).toContain('sd-ran-gui app is running!');
  });
});
